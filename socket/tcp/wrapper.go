package tcp

import (
	"bufio"
	"encoding/binary"
	"github.com/xlkness/netgo/event"
	"github.com/xlkness/netgo/socket"
	"github.com/xlkness/netgo/utils"
	"net"
	"sync/atomic"
	"time"
)

func NewSocketTcpListener(addr string, maxRecvMsgLen int32, maxWriteDeadline time.Duration,
	eventCb func(*event.Event, socket.Socket)) *SocketTcp {
	listener := &SocketTcp{}
	listener.addr = addr
	listener.eventCb = eventCb
	listener.maxRecvLen = maxRecvMsgLen
	listener.maxWriteDeadline = maxWriteDeadline
	return listener
}

func NewSocketTcpConnector(addr string, maxRecvLen int32, maxWriteDeadline time.Duration) *SocketTcp {
	connector := &SocketTcp{}
	connector.addr = addr
	connector.maxRecvLen = maxRecvLen
	connector.maxWriteDeadline = maxWriteDeadline
	connector.writeChann = make(chan []byte, 1024)
	return connector
}

func (listener *SocketTcp) StartListen() error {
	listenFd, err := net.Listen("tcp", listener.addr)
	if err != nil {
		return err
	}
	go func() {
		for !listener.getIsStop() {
			clientFd, err := listenFd.Accept()
			if err != nil {
				listener.setIsStop(1)
				break
			}
			client := listener.newSocketTcpConnector(clientFd)
			socketEvent := event.NewEvent(event.EventTypeSocketConnect, 0, nil)
			client.eventCb(socketEvent, client)
			listener.handleConnectorSession(client)
		}
	}()
	return nil
}

func (listener *SocketTcp) EndListen() {
	listener.setIsStop(1)
}

func (connector *SocketTcp) Connect() error {
	conn, err := net.Dial("tcp", connector.addr)
	if err != nil {
		return err
	}
	connector.conn = conn
	connector.readStream = bufio.NewReader(conn)
	go connector.handleConnectorWrite()
	return nil
}

func (connector *SocketTcp) Read() (uint32, []byte, error) {
	return utils.ReadTLVMsg(connector.readStream, connector.maxRecvLen)
}

func (connector *SocketTcp) Write(tag uint32, payload []byte) error {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint32(buf, tag)
	binary.LittleEndian.PutUint32(buf[4:], uint32(len(payload)))
	buf = append(buf, payload...)
	connector.writeChann <- buf
	return nil
}

func (connector *SocketTcp) Close() {
	connector.conn.Close()
	connector.setIsStop(1)
}

func (su *SocketTcp) GetLocalAddr() string {
	return su.conn.LocalAddr().String()
}

func (su *SocketTcp) GetRemoteAddr() string {
	return su.conn.RemoteAddr().String()
}

func (listener *SocketTcp) newSocketTcpConnector(conn net.Conn) *SocketTcp {
	connector := &SocketTcp{}
	connector.conn = conn
	connector.readStream = bufio.NewReader(conn)
	connector.maxRecvLen = listener.maxRecvLen
	connector.maxWriteDeadline = listener.maxWriteDeadline
	connector.eventCb = listener.eventCb
	connector.writeChann = make(chan []byte, 1024)
	return connector
}

func (listener *SocketTcp) handleConnectorSession(connector *SocketTcp) {
	go connector.handleConnectorRead()
	go connector.handleConnectorWrite()
}

func (connector *SocketTcp) handleConnectorRead() {
	for !connector.getIsStop() {
		tag, payload, err := utils.ReadTLVMsg(connector.readStream, connector.maxRecvLen)
		if err != nil {
			connector.setIsStop(1)
			socketEvent := event.NewEvent(event.EventTypeSocketClose, tag, payload)
			connector.eventCb(socketEvent, connector)
			break
		}

		socketEvent := event.NewEvent(event.EventTypeSocketRecv, tag, payload)
		connector.eventCb(socketEvent, connector)
	}
}

func (connector *SocketTcp) handleConnectorWrite() {
	for !connector.getIsStop() {
		msg := <-connector.writeChann
		connector.conn.SetWriteDeadline(time.Now().Add(connector.maxWriteDeadline))
		_, err := connector.conn.Write(msg)
		if err != nil {
			connector.setIsStop(1)
			socketEvent := event.NewEvent(event.EventTypeSocketClose, 0, nil)
			connector.eventCb(socketEvent, connector)
			return
		}
	}
}

func (st *SocketTcp) setIsStop(val int32) {
	atomic.StoreInt32(&st.isStop, val)
}

func (st *SocketTcp) getIsStop() bool {
	return atomic.LoadInt32(&st.isStop) == 1
}

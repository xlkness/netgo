package udp

import (
	"encoding/binary"
	"fmt"
	"github.com/xlkness/netgo/event"
	"github.com/xlkness/netgo/netgo"
	"net"
	"sync/atomic"
	"time"
)

func NewSocketUdpListener(addr string, maxRecvMsgLen int32, maxWriteDeadline time.Duration,
	eventCb func(*event.Event, netgo.Socket)) *SocketUdp {
	listener := &SocketUdp{
		addr:       addr,
		maxRecvLen: maxRecvMsgLen,
	}
	listener.eventCb = eventCb
	return listener
}

func NewSocketUdpConnector(addr string, maxRecvMsgLen int32, maxWriteDeadline time.Duration) *SocketUdp {
	connector := &SocketUdp{
		addr: addr,
	}
	connector.maxRecvLen = maxRecvMsgLen
	return connector
}

func (listener *SocketUdp) StartListen() error {
	udpAddr, err := net.ResolveUDPAddr("udp", listener.addr)
	if err != nil {
		return err
	}
	listenerFd, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	listener.conn = listenerFd
	atomic.StoreInt32(&listener.isListener, 1)
	buf := make([]byte, listener.maxRecvLen)
	for !listener.getIsStop() {

		n, cAddr, err := listenerFd.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("read from udp err:%v\n", err)
			listener.setIsStop(1)
			return err
		}
		if n < 8 {
			listener.setIsStop(1)
			return fmt.Errorf("recv msg length < 8:%v", n)
		}
		tag := binary.LittleEndian.Uint32(buf)
		length := binary.LittleEndian.Uint32(buf[4:])
		payload := buf[8 : length+8]
		e := event.NewEvent(event.EventTypeSocketRecv, tag, payload)
		client := listener.newSocketUdpConnector(listenerFd, cAddr)
		listener.eventCb(e, client)
	}

	return nil
}

func (listener *SocketUdp) EndListen() {
	listener.setIsStop(1)
	listener.conn.Close()
}

func (connector *SocketUdp) Connect() error {
	conn, err := net.Dial("udp", connector.addr)
	if err != nil {
		return err
	}
	connector.conn = conn
	return nil
}

// only use for connector
func (connector *SocketUdp) Read() (uint32, []byte, error) {
	buf := make([]byte, connector.maxRecvLen)
	n, err := connector.conn.Read(buf)
	if err != nil {
		return 0, nil, err
	}
	if n < 8 {
		return 0, nil, fmt.Errorf("recv msg length < 8:%v", n)
	}

	tag := binary.LittleEndian.Uint32(buf)
	length := binary.LittleEndian.Uint32(buf[4:])

	if length > uint32(connector.maxRecvLen) {
		return 0, nil, fmt.Errorf("recv msg max length:%v/%v", length, connector.maxRecvLen)
	}
	if length != uint32(n-8) {
		return 0, nil, fmt.Errorf("recv msg invalid length:%v/%v", length, n)
	}
	payload := buf[8 : length+8]
	return tag, payload, nil
}

func (su *SocketUdp) Write(tag uint32, payload []byte) error {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint32(buf, tag)
	binary.LittleEndian.PutUint32(buf[4:], uint32(len(payload)))
	buf = append(buf, payload...)
	if atomic.LoadInt32(&su.isListener) == 1 {
		_, err := su.conn.(*net.UDPConn).WriteToUDP(buf, su.udpAddr)
		if err != nil {
			return err
		}
	} else {
		_, err := su.conn.Write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (connector *SocketUdp) Close() {
	connector.conn.Close()
}

func (su *SocketUdp) GetLocalAddr() string {
	return su.conn.LocalAddr().String()
}

func (su *SocketUdp) GetRemoteAddr() string {
	if atomic.LoadInt32(&su.isListener) == 1 {
		return su.udpAddr.String()
	}
	return su.conn.RemoteAddr().String()
}

func (listener *SocketUdp) newSocketUdpConnector(conn net.Conn, cAddr *net.UDPAddr) *SocketUdp {
	connector := &SocketUdp{}
	connector.conn = conn
	connector.udpAddr = cAddr
	connector.maxRecvLen = listener.maxRecvLen
	connector.eventCb = listener.eventCb
	atomic.StoreInt32(&connector.isListener, 1)
	return connector
}

func (st *SocketUdp) setIsStop(val int32) {
	atomic.StoreInt32(&st.isStop, val)
}

func (st *SocketUdp) getIsStop() bool {
	return atomic.LoadInt32(&st.isStop) == 1
}

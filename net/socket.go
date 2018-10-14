package net

import (
	"github.com/xlkness/netgo/event"
	"github.com/xlkness/netgo/net/tcp"
	"github.com/xlkness/netgo/net/udp"
	"net"
	"time"
)

type Socket interface {
	StartListen() error
	EndListen()
	Connect() error
	// only use for connector
	Read() (uint32, []byte, error)
	Write(msg []byte) error
	Close()
	GetRawConn() net.Conn
}

func NewSocketListener(commType, addr string, maxRecvMsgLen int32, maxWriteDeadline time.Duration,
	eventcb func(*event.Event, Socket)) Socket {
	if commType == "tcp" {
		return tcp.NewSocketTcpListener(addr, maxRecvMsgLen, maxWriteDeadline, eventcb)
	} else if commType == "udp" {
		return udp.NewSocketUdpListener(addr, maxRecvMsgLen, maxWriteDeadline, eventcb)
	}
	return nil
}

func NewSocketConnector(commType, addr string, maxRecvMsgLen int32, maxWriteDeadline time.Duration) Socket {
	if commType == "tcp" {
		return tcp.NewSocketTcpConnector(addr, maxRecvMsgLen, maxWriteDeadline)
	} else if commType == "udp" {
		return udp.NewSocketUdpConnector(addr, maxRecvMsgLen, maxWriteDeadline)
	}
	return nil
}

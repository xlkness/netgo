package netgo

import (
	"github.com/xlkness/netgo/event"
	"github.com/xlkness/netgo/socket"
	"github.com/xlkness/netgo/socket/tcp"
	"github.com/xlkness/netgo/socket/udp"
	"time"
)

func NewSocketListener(commType, addr string, maxRecvMsgLen int32, maxWriteDeadline time.Duration,
	eventcb func(*event.Event, socket.Socket)) socket.Socket {
	if commType == "tcp" {
		return tcp.NewSocketTcpListener(addr, maxRecvMsgLen, maxWriteDeadline, eventcb)
	} else if commType == "udp" {
		return udp.NewSocketUdpListener(addr, maxRecvMsgLen, maxWriteDeadline, eventcb)
	}
	return nil
}

func NewSocketConnector(commType, addr string, maxRecvMsgLen int32, maxWriteDeadline time.Duration) socket.Socket {
	if commType == "tcp" {
		return tcp.NewSocketTcpConnector(addr, maxRecvMsgLen, maxWriteDeadline)
	} else if commType == "udp" {
		return udp.NewSocketUdpConnector(addr, maxRecvMsgLen, maxWriteDeadline)
	}
	return nil
}

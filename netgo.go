package netgo

import (
	"github.com/xlkness/netgo/event"
	"github.com/xlkness/netgo/netgo"
	"github.com/xlkness/netgo/netgo/tcp"
	"github.com/xlkness/netgo/netgo/udp"
	"time"
)

func NewSocketListener(commType, addr string, maxRecvMsgLen int32, maxWriteDeadline time.Duration,
	eventcb func(*event.Event, netgo.Socket)) netgo.Socket {
	if commType == "tcp" {
		return tcp.NewSocketTcpListener(addr, maxRecvMsgLen, maxWriteDeadline, eventcb)
	} else if commType == "udp" {
		return udp.NewSocketUdpListener(addr, maxRecvMsgLen, maxWriteDeadline, eventcb)
	}
	return nil
}

func NewSocketConnector(commType, addr string, maxRecvMsgLen int32, maxWriteDeadline time.Duration) netgo.Socket {
	if commType == "tcp" {
		return tcp.NewSocketTcpConnector(addr, maxRecvMsgLen, maxWriteDeadline)
	} else if commType == "udp" {
		return udp.NewSocketUdpConnector(addr, maxRecvMsgLen, maxWriteDeadline)
	}
	return nil
}

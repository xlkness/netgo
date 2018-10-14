package udp

import (
	"github.com/xlkness/netgo/event"
	net1 "github.com/xlkness/netgo/netgo"
	"net"
	"time"
)

type Listener struct {
	eventCb func(*event.Event, net1.Socket)
}

type Connector struct {
}

type SocketUdp struct {
	Listener
	Connector
	isListener       int32
	addr             string
	conn             net.Conn
	udpAddr          *net.UDPAddr
	isStop           int32
	maxRecvLen       int32
	maxWriteDeadline time.Duration
}

package udp

import (
	"github.com/xlkness/netgo/event"
	"net"
	"time"
)

type Listener struct {
	eventCb func(*event.Event, interface{})
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

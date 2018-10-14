package tcp

import (
	"bufio"
	"github.com/xlkness/netgo/event"
	"net"
	"time"
)

type Listener struct {
	eventCb func(*event.Event, interface{})
}

type Connector struct {
	writeChann chan []byte
}

type SocketTcp struct {
	Listener
	Connector
	addr             string
	conn             net.Conn
	isStop           int32
	readStream       *bufio.Reader
	maxRecvLen       int32
	maxWriteDeadline time.Duration
}

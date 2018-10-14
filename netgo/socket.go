package netgo

import "net"

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

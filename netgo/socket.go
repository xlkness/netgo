package netgo

type Socket interface {
	StartListen() error
	EndListen()
	Connect() error
	// only use for connector
	Read() (uint32, []byte, error)
	Write(msg []byte) error
	Close()
	GetLocalAddr() string
	GetRemoteAddr() string
}

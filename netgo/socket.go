package netgo

type Socket interface {
	StartListen() error
	EndListen()
	Connect() error
	// only use for connector
	Read() (uint32, []byte, error)

	/*
		memory layout in a frame:
		|--------------------------------------------------------|
		|--tag:4bytes--|--length:4bytes--|--payload:maxLength-8--|
		|--------------------------------------------------------|
		note: tag/length are both using little endian
	*/
	Write(tag uint32, payload []byte) error
	Close()
	GetLocalAddr() string
	GetRemoteAddr() string
}

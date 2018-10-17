package tcp

import (
	"fmt"
	"github.com/xlkness/netgo/event"
	"github.com/xlkness/netgo/socket"
	"sync"
	"testing"
	"time"
)

func TestEcho(t *testing.T) {
	server := NewSocketTcpListener("127.0.0.1:9190", 1<<16, time.Minute, func(e *event.Event, client socket.Socket) {
		switch e.Type {
		case event.EventTypeSocketConnect:
			fmt.Printf("client connect!\n")
		case event.EventTypeSocketClose:
			fmt.Printf("client close!\n")
		case event.EventTypeSocketRecv:
			fmt.Printf("recv client msg:%v-%v\n", e.Msg.Tag, string(e.Msg.Payload))
			client.Write(e.Msg.Tag, e.Msg.Payload)
		}
	})
	server.StartListen()

	wg := &sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			c := NewSocketTcpConnector("127.0.0.1:9190", 1<<16, time.Minute)
			err := c.Connect()
			if err != nil {
				fmt.Printf("connect %s err:%v\n", c.addr, err)
				return
			}
			c.Write(uint32(123), []byte("netgo"))
			tag, msg, err := c.Read()
			if err != nil {
				fmt.Printf("write and recv err:%v", err)
				return
			}
			fmt.Printf("recv echo msg:%v-%v\n", tag, string(msg))
			c.Close()
			wg.Done()
		}()
	}
	wg.Wait()
}

package udp

import (
	"encoding/binary"
	"fmt"
	"github.com/xlkness/netgo/event"
	"github.com/xlkness/netgo/socket"
	"sync"
	"testing"
	"time"
)

func TestEcho(t *testing.T) {
	wg := &sync.WaitGroup{}
	server := NewSocketUdpListener("127.0.0.1:9091", 1<<16, -1, func(e *event.Event, client socket.Socket) {
		switch e.Type {
		case event.EventTypeSocketRecv:
			fmt.Printf("recv msg:%v\n", string(e.Msg.Payload))
			err := client.Write(e.Msg.Tag, e.Msg.Payload)
			if err != nil {
				fmt.Printf("echo write err:%v, %v\n", err, len(e.Msg.Payload))
			}
		}
	})
	go server.StartListen()
	time.Sleep(1 * time.Second)
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			client := NewSocketUdpConnector("127.0.0.1:9091", 1<<16, -1)
			if err := client.Connect(); err != nil {
				fmt.Printf("conect err:%v\n", err)
				return
			}
			err := client.Write(uint32(123), []byte("netgo"))
			if err != nil {
				fmt.Printf("write err:%v\n", err)
				return
			}
			tag, payload, err := client.Read()
			if err != nil {
				fmt.Printf("recv err:%v\n", err)
				return
			}
			fmt.Printf("recv echo msg:%v-%v\n", tag, string(payload))
			wg.Done()
		}()
	}
	wg.Wait()
}

package udp

import (
	"encoding/binary"
	"fmt"
	"github.com/xlkness/netgo/event"
	"github.com/xlkness/netgo/netgo"
	"sync"
	"testing"
	"time"
)

func TestEcho(t *testing.T) {
	wg := &sync.WaitGroup{}
	server := NewSocketUdpListener("127.0.0.1:9091", 1<<16, -1, func(e *event.Event, client netgo.Socket) {
		switch e.Type {
		case event.EventTypeSocketRecv:
			fmt.Printf("recv msg:%v\n", string(e.Msg.Payload))
			buf := make([]byte, 8)
			binary.LittleEndian.PutUint32(buf, e.Msg.Tag)
			binary.LittleEndian.PutUint32(buf[4:], uint32(len(e.Msg.Payload)))
			buf = append(buf, e.Msg.Payload...)
			err := client.Write(buf)
			if err != nil {
				fmt.Printf("echo write err:%v, %v\n", err, len(buf))
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
			payload := []byte("likun")
			buf := make([]byte, 8)
			binary.LittleEndian.PutUint32(buf, uint32(123))
			binary.LittleEndian.PutUint32(buf[4:], uint32(len(payload)))
			buf = append(buf, payload...)
			err := client.Write(buf)
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

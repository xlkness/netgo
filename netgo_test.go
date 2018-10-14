package netgo

import (
	"encoding/binary"
	"fmt"
	"github.com/xlkness/netgo/event"
	"github.com/xlkness/netgo/netgo"
	"sync"
	"testing"
	"time"
)

func TestEchoTcp(t *testing.T) {
	testLogicFun(t, "tcp")
}

func TestEchoUdp(t *testing.T) {
	testLogicFun(t, "udp")
}

func testLogicFun(t *testing.T, commType string) {
	wg := &sync.WaitGroup{}

	server := NewSocketListener(commType, "192.168.1.188:9190", 1<<16, time.Minute, func(e *event.Event, client netgo.Socket) {
		switch e.Type {
		case event.EventTypeSocketConnect:
			fmt.Printf("%v <- %v, client connect!\n", client.GetLocalAddr(), client.GetRemoteAddr())
		case event.EventTypeSocketClose:
			fmt.Printf("%v <- %v, client closed!\n", client.GetLocalAddr(), client.GetRemoteAddr())
		case event.EventTypeSocketRecv:
			fmt.Printf("%v <- %v, client packet:%v-%v\n", client.GetLocalAddr(), client.GetRemoteAddr(), e.Msg.Tag, string(e.Msg.Payload))
			client.Write(getTestPacket(e.Msg.Tag, e.Msg.Payload))
		}
	})

	go server.StartListen()

	// delay for complete start udp server
	time.Sleep(1 * time.Second)

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			client := NewSocketConnector(commType, "192.168.1.188:9190", 1<<16, time.Minute)
			err := client.Connect()
			if err != nil {
				t.Fatalf("client connect err:%v\n", err)
				return
			}
			err = client.Write(getTestPacket(uint32(12345), []byte("netgo")))
			if err != nil {
				t.Fatalf("client write err:%v\n", err)
				return
			}
			tag, payload, err := client.Read()
			if err != nil {
				t.Fatalf("client read err:%v\n", err)
				return
			}
			fmt.Printf("client recv echo msg:%v-%v\n", tag, string(payload))
			client.Close()
			wg.Done()
		}()
	}
	wg.Wait()
}

func getTestPacket(tag uint32, payload []byte) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint32(buf, tag)
	binary.LittleEndian.PutUint32(buf[4:], uint32(len(payload)))
	buf = append(buf, payload...)
	return buf
}

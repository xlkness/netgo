package net

import (
	"encoding/binary"
	"fmt"
	"github.com/xlkness/netgo/event"
	"testing"
	"time"
)

func TestEchoTcp(t *testing.T) {
	server := NewSocketListener("tcp", "192.168.1.188:9190", 1<<16, time.Minute, func(e *event.Event, client Socket) {
		switch e.Type {
		case event.EventTypeSocketConnect:
			fmt.Printf("client:%v connect!\n", client.GetRawConn().RemoteAddr().String())
		case event.EventTypeSocketClose:
			fmt.Printf("client:%v closed!\n", client.GetRawConn().RemoteAddr().String())
		case event.EventTypeSocketRecv:
			fmt.Printf("client:%v packet:%v-%v\n", client.GetRawConn().RemoteAddr().String(), e.Msg.Tag, string(e.Msg.Payload))
			client.Write(getTestPacket(e.Msg.Tag, e.Msg.Payload))
		}
	})
	server.StartListen()
	for i := 0; i < 10; i++ {
		go func() {
			client := NewSocketConnector("tcp", "192.168.1.188:9190", 1<<16, time.Minute)
			err := client.Connect()
			if err != nil {
				fmt.Printf("client connect err:%v\n", err)
				return
			}
			err = client.Write(getTestPacket(uint32(12345), []byte("netgo")))
			if err != nil {
				fmt.Printf("client write err:%v\n", err)
				return
			}
			tag, payload, err := client.Read()
			if err != nil {
				fmt.Printf("client read err:%v\n", err)
				return
			}
			fmt.Printf("client recv echo msg:%v-%v\n", tag, string(payload))
		}()
	}
}

func getTestPacket(tag uint32, payload []byte) []byte {

	buf := make([]byte, 8)
	binary.LittleEndian.PutUint32(buf, tag)
	binary.LittleEndian.PutUint32(buf, uint32(len(payload)))
	buf = append(buf, payload...)
	return buf
}

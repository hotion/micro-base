package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/shiguanghuxian/micro-base/kit/transport/tcppacket"
	"github.com/shiguanghuxian/tcplibrary"
)

/* tcp消息测试 */

func tcpHello() {
	for i := 0; i < 100; i++ {
		client := new(Client)
		c, err := tcplibrary.NewTCPClient(true, client, &tcppacket.MicroPacket{})
		if err != nil {
			log.Println(err)
		}
		err = c.DialAndStart(":38080")
		log.Println(err)
	}
}

// Client tcp客户端服务对象
type Client struct {
}

// OnConnect 连接建立时
func (c *Client) OnConnect(conn *tcplibrary.Conn) error {
	log.Println("OnConnect")
	go func() {
		conn := conn
		i := 1
		for {
			pp := &tcppacket.MicroPacket{
				Payload:      `{"name":"测试TCP"}`,
				EndpointType: tcppacket.TCPPostHelloEndpoint,
				Sequence:     uint16(i),
			}
			n, err := conn.SendMessage(pp)
			log.Println(n, err)
			i++
			if i > 60000 {
				i = 1
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return nil
}

// OnError 遇到错误时
func (c *Client) OnError(err error) {
	log.Println("OnError")
	log.Println(err)
}

// OnClose 连接关闭时
func (c *Client) OnClose(conn *tcplibrary.Conn, err error) {
	log.Println("OnClose")
	log.Println(err)
	os.Exit(1)
}

// OnRecMessage 收到消息时
func (c *Client) OnRecMessage(conn *tcplibrary.Conn, v interface{}) {
	log.Println("OnRecMessage")
	log.Println(v)
	if packet, ok := v.(*tcppacket.MicroPacket); ok == true {
		log.Printf("消息体长度:%d 消息体内容:%s", packet.Length, string(packet.Payload))
	} else {
		js, _ := json.Marshal(v)
		log.Println(string(js))
	}
}

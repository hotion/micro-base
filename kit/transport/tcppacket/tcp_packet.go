package tcppacket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"golang.org/x/net/websocket"
)

/*
tcp 方式包打包和解包对象
*/

// MicroPacket 实现自定义协议，包的解析和编码
type MicroPacket struct {
	Header         string           `json:"header"`        // 包头 固定为 `VI` [0x56,0x49]
	Length         uint16           `json:"length"`        // 包长度,2字节
	EndpointType   TCPEndpointType  `json:"endpoint_type"` // 消息对应的端点
	Sequence       uint16           `json:"sequence"`      // 报文序号,客户端自增，服务端返回永远为0
	Reserve        int32            `json:"reserve"`       // 预留
	Payload        string           `json:"payload"`       // 报文内容,n字节
	WebsocketCodec *websocket.Codec `json:"-"`
}

// TCPEndpointType tcp 端点类型，每个kit端点对应一个枚举值
type TCPEndpointType int16

const (
	// TCPPostHelloEndpoint tcp hello 端点
	TCPPostHelloEndpoint TCPEndpointType = iota // 0
	// ...
)

// Unmarshal 默认解包
func (mp *MicroPacket) Unmarshal(data []byte, c chan interface{}) ([]byte, error) {
	var err error
	// 捕获异常
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%T", r)
		}
	}()
	// 长度不足4个字节无法获取包长度
	if len(data) < 12 {
		return data, err
	}
	// 截取前两个字节判断是否是 'VI'
	packetHeader := data[:2]
	if !bytes.Equal(packetHeader, []byte{0x56, 0x49}) {
		return data, errors.New("解包错误，包头必须是VI")
	}
	// 获取包长度
	packetLength := BytesToInt16(data[2:4])
	packetLength = packetLength + 13
	// 判断data是否大于一个报文长度
	if len(data) < int(packetLength) {
		return data, err
	}
	// 截取一个包的长度，解包
	packetData := data[:packetLength]
	// 调用单独解包函数
	packet, err := mp.UnmarshalOne(packetData)
	if err != nil {
		return mp.Unmarshal(data[packetLength:], c)
	}
	// 向消息管道写消息
	c <- packet
	// 重新调用自己，拆包
	return mp.Unmarshal(data[packetLength:], c)
}

// UnmarshalOne 将一个包长的字节转为一个对象
func (mp *MicroPacket) UnmarshalOne(data []byte) (*MicroPacket, error) {
	if len(data) < 13 {
		return nil, errors.New("包长不足13个字节")
	}
	packet := new(MicroPacket)
	packet.Header = "VI"
	packet.Length = BytesToUint16(data[2:4])
	packet.EndpointType = TCPEndpointType(BytesToInt16(data[4:6]))
	packet.Sequence = BytesToUint16(data[6:8])
	packet.Reserve = BytesToInt32(data[8:12])
	packet.Payload = string(data[12:])
	return packet, nil
}

// Marshal 默认封包
func (mp *MicroPacket) Marshal(v interface{}) ([]byte, error) {
	packet, ok := v.(*MicroPacket)
	if ok == false {
		return nil, errors.New("封包参数不是*MicroPacket")
	}
	// 创建Buffer对象，并写入头 'VI'
	packetData := bytes.NewBuffer([]byte{})
	packetData.Write([]byte{0x56, 0x49})
	// 计算包长度
	packetLengthBytes := IntToBytes(int16(len([]byte(packet.Payload))))
	packetData.Write(packetLengthBytes)
	// 消息对应的端点
	packetData.Write(IntToBytes(int16(packet.EndpointType)))
	// 报文序号
	packetData.Write(IntToBytes(packet.Sequence))
	// 写入预留
	packetData.Write(IntToBytes(packet.Reserve))

	// 判断包头长度是否是 12
	if packetData.Len() != 12 {
		return nil, errors.New("报文长度不是12")
	}
	// 写入内容
	packetData.WriteString(packet.Payload)

	return packetData.Bytes(), nil
}

// MarshalToJSON 编码到json, 同时将Payload转为字符串
func (mp *MicroPacket) MarshalToJSON(v interface{}) (data []byte, payloadType byte, err error) {
	packet, ok := v.(*MicroPacket)
	if ok == false {
		return []byte{}, websocket.TextFrame, errors.New("封包参数不是*MicroPacket")
	}
	data, err = json.Marshal(packet)
	return data, websocket.TextFrame, err
}

// UnmarshalToJSON 解包为json字符串形式
func (mp *MicroPacket) UnmarshalToJSON(data []byte, payloadType byte, v interface{}) (err error) {
	return json.Unmarshal(data, v)
}

// GetWebsocketCodec 获取websocket编解码对象
func (mp *MicroPacket) GetWebsocketCodec() *websocket.Codec {
	return mp.WebsocketCodec
}

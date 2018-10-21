package transport

import (
	"context"
	"encoding/json"

	uuid "github.com/satori/go.uuid"
	"github.com/shiguanghuxian/micro-base/internal/log"
	"github.com/shiguanghuxian/micro-base/kit/endpoint"
	"github.com/shiguanghuxian/micro-base/kit/service"
	"github.com/shiguanghuxian/micro-base/kit/transport/tcppacket"
	"github.com/shiguanghuxian/tcplibrary"
)

/*
 tcp方式对外提供服务
 使用 https://github.com/shiguanghuxian/tcplibrary
*/

type tcpServer struct {
	server    *tcplibrary.TCPServer // 服务监听对象
	endpoints endpoint.Endpoints
	service   service.Service
	logger    *log.Log
}

// NewTCPHandler 创建tcp服务
func NewTCPHandler(endpoints endpoint.Endpoints, s service.Service, logger *log.Log) (*tcpServer, error) {
	srv := &tcpServer{
		endpoints: endpoints,
		service:   s,
		logger:    logger,
	}
	liSrv, err := tcplibrary.NewTCPServer(true, srv)
	if err != nil {
		return nil, err
	}
	srv.server = liSrv

	return srv, nil
}

// 启动服务
func (s *tcpServer) ListenAndServe(address string) error {
	return s.server.ListenAndServe(address)
}

// OnConnect 连接建立时
func (s *tcpServer) OnConnect(conn *tcplibrary.Conn) error {
	s.logger.Infow("有客户端建立连接")
	return nil
}

// OnError 连接遇到错误时
func (s *tcpServer) OnError(err error) {
	s.logger.Errorw("连接出现错误", "err", err)
}

// OnClose 连接关闭时
func (s *tcpServer) OnClose(conn *tcplibrary.Conn, err error) {
	if err != nil {
		s.logger.Errorw("关闭连接出现错误", "err", err)
	}
}

// OnRecMessage 收到客户端发送过来的消息时
func (s *tcpServer) OnRecMessage(conn *tcplibrary.Conn, v interface{}) {
	if packet, ok := v.(*tcppacket.MicroPacket); ok == true {
		s.SelectEndpoint(packet)
	} else {
		s.logger.Errorw("收到消息不是MicroPacket对象", "data", v)
	}
}

// GetClientID 生成一个客户端连接，只要唯一即可
func (s *tcpServer) GetClientID() string {
	newUUID, _ := uuid.NewV4()
	return newUUID.String()
}

// 根据消息端点类型，调用制定函数
func (s *tcpServer) SelectEndpoint(packet *tcppacket.MicroPacket) {
	switch packet.EndpointType {
	case tcppacket.TCPPostHelloEndpoint:
		req := new(endpoint.PostHelloRequest)
		err := json.Unmarshal([]byte(packet.Payload), req)
		if err == nil {
			s.endpoints.PostHello(context.Background(), req.Name)
		} else {
			s.logger.Errorw("json解析为PostHelloRequest错误", "err", err)
		}
		break
	default:
		s.logger.Warnw("未找到消息中的端点", "endpoint", packet.EndpointType, "packet", packet)
	}
}

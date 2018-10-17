package transport

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/shiguanghuxian/micro-base/internal/log"
	"github.com/shiguanghuxian/micro-base/kit/endpoint"
	"github.com/shiguanghuxian/micro-base/kit/service"
	"github.com/shiguanghuxian/micro-base/pb"
	"go.uber.org/zap"
)

type grpcServer struct {
	postHello grpctransport.Handler
}

// NewGRPCServer 创建grpc服务
func NewGRPCServer(endpoints endpoint.Endpoints, s service.Service, logger *zap.SugaredLogger) pb.HelloServer {
	// 配置
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(log.NewLog(logger)),
	}

	return &grpcServer{
		postHello: grpctransport.NewServer(
			endpoints.PostHelloEndpoint,
			decodeGRPCPostHelloRequest,
			encodeGRPCPostHelloResponse,
			options...,
		),
	}
}

// 此函数会调用 endpoint层 进而调用服务层
func (s grpcServer) PostHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	_, rep, err := s.postHello.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.HelloReply), nil
}

// Server 处理grpc请求参数解析，将 grpc层请求类型 转换为 endpoint层 数据类型
func decodeGRPCPostHelloRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.HelloRequest)
	return endpoint.PostHelloRequest{Name: req.GetName()}, nil
}

// Server 处理grpc响应参数解析，将 endpoint层响应类型 参数类型转换为grpc响应类型
func encodeGRPCPostHelloResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoint.PostHelloResponse)
	return &pb.HelloReply{Word: resp.Word}, nil
}

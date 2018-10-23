package transport

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/shiguanghuxian/micro-base/kit/endpoint"
	"github.com/shiguanghuxian/micro-base/model"
	"github.com/shiguanghuxian/micro-base/pb"
	"github.com/shiguanghuxian/micro-common/log"
)

/* grpc方式对外提供服务 */

type grpcServer struct {
	login grpctransport.Handler
}

// NewGRPCServer 创建grpc服务
func NewGRPCServer(endpoints endpoint.Endpoints, logger *log.Log) pb.AccountServer {
	// 配置
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}

	return &grpcServer{
		login: grpctransport.NewServer(
			endpoints.LoginEndpoint,
			decodeGRPCLoginRequest,
			encodeGRPCLoginResponse,
			options...,
		),
	}
}

// 登录
func (s grpcServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	_, rep, err := s.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.LoginReply), nil
}

func decodeGRPCLoginRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.LoginRequest)
	return &endpoint.LoginRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}, nil
}

func encodeGRPCLoginResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*model.User)
	return &pb.LoginReply{
		Username: resp.Username,
		Nickname: resp.Nickname,
		Password: resp.Password,
		Token:    resp.Token,
	}, nil
}

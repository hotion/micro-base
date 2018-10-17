package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/shiguanghuxian/micro-base/kit/service"
)

// MakeHelloEndpoint demo
func MakeHelloEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PostHelloRequest)
		str, err := s.PostHello(ctx, req.Name)
		return PostHelloResponse{Word: str}, err
	}
}

// PostHelloRequest 请求参数
type PostHelloRequest struct {
	Name string `json:"name"`
}

// PostHelloResponse 响应数据
type PostHelloResponse struct {
	Word string `json:"word"`
}

// PostHello 客户端使用到了
func (s *Endpoints) PostHello(ctx context.Context, name string) (string, error) {
	resp, err := s.PostHelloEndpoint(ctx, PostHelloRequest{Name: name})
	if err != nil {
		return "", err
	}
	response := resp.(PostHelloResponse)
	return response.Word, err
}

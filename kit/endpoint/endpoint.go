package endpoint

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/shiguanghuxian/micro-base/kit/service"
)

// Endpoints 将所有端点收集到一个结构体中
type Endpoints struct {
	PostHelloEndpoint endpoint.Endpoint
}

// MakeServerEndpoints 创建端点
func MakeServerEndpoints(s service.Service) Endpoints {
	return Endpoints{
		PostHelloEndpoint: MakeHelloEndpoint(s),
	}
}

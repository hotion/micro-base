package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/shiguanghuxian/micro-base/kit/service"
)

/* 对外提供服务端点 */

// MakeLoginEndpoint 用户登录
func MakeLoginEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*LoginRequest)
		user, err := s.Login(ctx, req.Username, req.Password)
		return user, FormatError(err)
	}
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

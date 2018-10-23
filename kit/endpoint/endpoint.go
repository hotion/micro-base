package endpoint

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/shiguanghuxian/micro-base/kit/service"
	"github.com/shiguanghuxian/micro-common/microerror"
)

// Endpoints 将所有端点收集到一个结构体中
type Endpoints struct {
	LoginEndpoint endpoint.Endpoint
}

// MakeServerEndpoints 创建端点
func MakeServerEndpoints(s service.Service) Endpoints {
	return Endpoints{
		LoginEndpoint: MakeLoginEndpoint(s),
	}
}

// FormatError 整理返回的错误，统一为*microerror.MicroError
func FormatError(err error) (microError error) {
	if err == nil {
		return nil
	}
	microError, ok := err.(*microerror.MicroError)
	if ok == false {
		microError = microerror.GetMicroError(30000, err)
	}
	return
}

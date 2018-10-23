package transport

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/shiguanghuxian/micro-base/kit/endpoint"
	"github.com/shiguanghuxian/micro-base/model"
	"github.com/shiguanghuxian/micro-common/log"
)

/* http方式对外提供服务 */

// NewHTTPHandler http服务定义
func NewHTTPHandler(endpoints endpoint.Endpoints, logger *log.Log) http.Handler {
	r := mux.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
	}

	// 路由
	r.Methods("POST").Path("/login").Handler(httptransport.NewServer(
		endpoints.LoginEndpoint,
		decodeLoginRequest,
		encodeLoginResponse,
		options...,
	))

	return r
}

func decodeLoginRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	req := new(endpoint.LoginRequest)
	if e := json.NewDecoder(r.Body).Decode(req); e != nil {
		return nil, e
	}
	return req, nil
}

func encodeLoginResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if _, ok := response.(*model.User); ok == false {
		// 返回错误码 - 应该输出错误
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

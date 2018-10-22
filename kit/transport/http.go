package transport

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/shiguanghuxian/micro-base/kit/endpoint"
	"github.com/shiguanghuxian/micro-base/kit/service"
	"github.com/shiguanghuxian/micro-common/log"
)

/* http方式对外提供服务 */

// NewHTTPHandler http服务定义
func NewHTTPHandler(endpoints endpoint.Endpoints, s service.Service, logger *log.Log) http.Handler {
	r := mux.NewRouter()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
	}

	// POST Hello Word
	r.Methods("POST").Path("/hello").Handler(httptransport.NewServer(
		endpoints.PostHelloEndpoint,
		decodeHelloRequest,
		encodeHelloResponse,
		options...,
	))

	return r
}

func decodeHelloRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req endpoint.PostHelloRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func encodeHelloResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	// fmt.Printf("response is %T", response)
	if _, ok := response.(endpoint.PostHelloResponse); ok == false {
		// 返回错误码 - 应该输出错误
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

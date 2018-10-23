package logging

import (
	"context"

	"github.com/shiguanghuxian/micro-base/kit/service"
	"github.com/shiguanghuxian/micro-base/model"
	"github.com/shiguanghuxian/micro-common/log"
)

// Middleware 中间件
type Middleware func(service.Service) service.Service

// LoggingMiddleware 注册日志中间件
func LoggingMiddleware(logger *log.Log) Middleware {
	return func(next service.Service) service.Service {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger *log.Log
	next   service.Service
}

// PostHello 服务层 Hello 日志
func (mw loggingMiddleware) Login(ctx context.Context, username, password string) (user *model.User, err error) {
	defer func() {
		if err != nil {
			mw.logger.Errorw("PostHello error", "method", "Login", "input username", username, "input password", password, "err", err)
		} else {
			mw.logger.Infow("PostHello info", "method", "Login", "input username", username, "input password", password)
		}
	}()
	user, err = mw.next.Login(ctx, username, password)
	return
}

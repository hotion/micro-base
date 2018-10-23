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
func (mw loggingMiddleware) PostHello(ctx context.Context, name string) (word string, err error) {
	defer func() {
		if err != nil {
			mw.logger.Errorw("PostHello error", "method", "PostHello", "input name", name, "output word", word)
		} else {
			mw.logger.Infow("PostHello info", "method", "PostHello", "input name", name, "output word", word)
		}
	}()
	word, err = mw.next.PostHello(ctx, name)
	return
}

// PostHello 服务层 Hello 日志
func (mw loggingMiddleware) Login(ctx context.Context, username, password string) (user *model.User, err error) {
	defer func() {
		if err != nil {
			mw.logger.Errorw("PostHello error", "method", "Login", "input username", username, "output password", password, "err", err)
		} else {
			mw.logger.Infow("PostHello info", "method", "Login", "input username", username, "output password", password)
		}
	}()
	user, err = mw.next.Login(ctx, username, password)
	return
}

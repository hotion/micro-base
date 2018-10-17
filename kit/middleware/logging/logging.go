package logging

import (
	"context"

	"github.com/shiguanghuxian/micro-base/kit/service"
	"go.uber.org/zap"
)

// Middleware 中间件
type Middleware func(service.Service) service.Service

// LoggingMiddleware 注册日志中间件
func LoggingMiddleware(logger *zap.SugaredLogger) Middleware {
	// logger1, _ := zap.NewProduction()
	// sugar := logger1.Sugar()

	return func(next service.Service) service.Service {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger *zap.SugaredLogger
	next   service.Service
}

// PostHello 服务层 Hello 日志
func (mw loggingMiddleware) PostHello(ctx context.Context, name string) (word string, err error) {
	defer func() {
		if err != nil {
			mw.logger.Error("method", "PostHello", "input name", name, "output word", word)
		} else {
			mw.logger.Info("method", "PostHello", "input name", name, "output word", word)
		}
	}()
	word, err = mw.next.PostHello(ctx, name)
	return
}

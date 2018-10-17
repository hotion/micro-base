package service

import (
	"context"
)

// Service 定义服务
type Service interface {
	PostHello(ctx context.Context, name string) (string, error)
}

// 基础service对象
type basicService struct {
}

// NewBasicService 创建基础服务对象
func NewBasicService() Service {
	return basicService{}
}

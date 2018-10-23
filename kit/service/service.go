package service

import (
	"context"

	"github.com/shiguanghuxian/micro-base/model"
)

// Service 定义服务
type Service interface {
	// 通过用户名密码登录
	Login(ctx context.Context, username, password string) (*model.User, error)
}

// 基础service对象
type basicService struct {
}

// NewBasicService 创建基础服务对象
func NewBasicService() Service {
	return basicService{}
}

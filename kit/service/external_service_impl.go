package service

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/shiguanghuxian/micro-base/model"
	"github.com/shiguanghuxian/micro-common/cache"
	"github.com/shiguanghuxian/micro-common/log"
	"github.com/shiguanghuxian/micro-common/microerror"
)

/* 对外提供服务实现 */

// Login 用户登录
func (s basicService) Login(ctx context.Context, username, password string) (user *model.User, err error) {
	user = new(model.User)
	err = user.UserLogin(username, password)
	if err != nil {
		log.Logger.Infow("登录操作数据库错误", "err", err)
		if err == gorm.ErrRecordNotFound {
			err = microerror.ErrRecordNotFound
		} else {
			err = microerror.ErrRecordNotFound
		}
		return
	}
	// 将登录信息缓存入redis
	key := cache.GetUserLoginToken()
	err = cache.GetClient().Set(key, user, 0).Err()
	if err != nil {

	}
}

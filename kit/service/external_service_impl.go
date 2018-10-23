package service

import (
	"context"
	"encoding/json"

	"github.com/jinzhu/gorm"
	"github.com/shiguanghuxian/micro-base/model"
	"github.com/shiguanghuxian/micro-common/cache"
	"github.com/shiguanghuxian/micro-common/crypto"
	"github.com/shiguanghuxian/micro-common/microerror"
)

/* 对外提供服务实现 */

// Login 用户登录
func (s basicService) Login(ctx context.Context, username, password string) (user *model.User, err error) {
	user = new(model.User)
	password = crypto.PasswordHash(password)
	err = user.UserLogin(username, password)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = microerror.GetMicroError(10001)
		}
		return
	}
	// 将登录信息缓存入redis
	key, token := cache.GetUserLoginToken()
	userJs, _ := json.Marshal(user)
	err = cache.GetClient().Set(key, string(userJs), 0).Err()
	if err != nil {
		err = microerror.GetMicroError(10102, err)
		return
	}

	user.Token = token
	user.Password = ""
	return
}

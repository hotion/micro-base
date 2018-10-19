package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"column:username"`
	Nickname string `json:"nickname" gorm:"column:nickname"`
	Password string `json:"password" gorm:"column:password"`
	IsDelete int8   `json:"is_delete" gorm:"column:is_delete"`
}

func (User) TableName() string {
	return gorm.DefaultTableNameHandler(nil, "user")
}

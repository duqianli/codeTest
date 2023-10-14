package models

import "gorm.io/gorm"

type UserBasic struct {
	gorm.Model
	Identity  string `gorm:"column:identity;type:varchar(36);"json:"identity"` //用户的唯一标识
	Name      string `gorm:"column:name;type:varchar(100);"json:"name"`
	Password  string `gorm:"column:password;type:varchar(32);"json:"password"`
	Phone     string `gorm:"column:phone;type:varchar(20);"json:"phone"`
	Mail      string `gorm:"column:mail;type:varchar(100);"json:"mail"`
	PassNum   int64  `gorm:"column:pass_num;type:int(11);" json:"pass_num"`     // 通过的次数
	SubmitNum int64  `gorm:"column:submit_num;type:int(11);" json:"submit_num"` // 提交次数
	IsAdmin   int    `gorm:"column:is_admin;type:tinyint(1);" json:"is_admin"`  // 是否是管理员【0-否，1-是】
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

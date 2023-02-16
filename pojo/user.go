package pojo

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName      string `json:"user_name" column:"user_name"` // 根据 username确定
	Password      string `json:"password"`
	FollowCount   int64  `json:"follow_count"`   // 关注数
	FollowerCount int64  `json:"follower_count"` // 粉丝数
}

func (User) TableName() string {
	return "users"
}

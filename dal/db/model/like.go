package model

import "gorm.io/gorm"

type Like struct {
	gorm.Model
	UserId  int64 `gorm:"column:user_id" json:"uid"`
	VideoId int64 `gorm:"column:video_id" json:"vid"`
	Action  bool  `gorm:"column:action" json:"action"`
}

func (Like) TableName() string {
	return "like"
}

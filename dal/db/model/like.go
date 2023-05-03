package model

import "gorm.io/gorm"

type Like struct {
	gorm.Model
	UserId  int64 `gorm:"column:user_id;uniqueIndex:idx_like_user_video" json:"uid"`
	VideoId int64 `gorm:"column:video_id;uniqueIndex:idx_like_user_video" json:"vid"`
	Action  bool  `gorm:"column:action" json:"action"`
}

func (Like) TableName() string {
	return "like"
}

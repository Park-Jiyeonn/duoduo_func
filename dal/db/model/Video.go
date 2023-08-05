package model

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	UserId        int    `gorm:"column:user_id;not null;index:fk_user_video" redis:"user_id"`
	Title         string `gorm:"column:title;type:varchar(128);not null" redis:"title"`
	PlayUrl       string `gorm:"column:play_url;varchar(128);not null" redis:"play_url"`
	CoverUrl      string `gorm:"column:cover_url;varchar(128);not null" redis:"cover_url"`
	FavoriteCount int    `gorm:"column:favorite_count;default:0" redis:"favorite_count"`
	CommentCount  int    `gorm:"column:comment_count;default:0" redis:"comment_count"`
}

func (v *Video) TableName() string {
	return "video"
}

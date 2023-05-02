package model

import (
	"gorm.io/gorm"
)

const TableNameComment = "comment"

type Comment struct {
	gorm.Model
	UserId      int64  `json:"user_id" gorm:"index:idx_comment_user_id"`
	VideoId     int64  `json:"video_id"`
	Content     string `json:"content"`
	PublishDate string `json:"publish_date"`
}

// TableName Comment's table name
func (*Comment) TableName() string {
	return TableNameComment
}

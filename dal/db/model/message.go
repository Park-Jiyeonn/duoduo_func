package model

import (
	"gorm.io/gorm"
)

const TableNameMessage = "message"

type Message struct {
	gorm.Model
	UserId      int64  `json:"user_info_id" gorm:"index:idx_message_user_id"`
	ToUserId    int64  `json:"video_id"`
	Content     string `json:"content"`
	PublishDate int64  `json:"publish_date"`
}

// TableName Message's table name
func (*Message) TableName() string {
	return TableNameMessage
}

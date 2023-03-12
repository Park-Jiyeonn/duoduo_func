package db

import (
	"context"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	UserID      int64  `json:"user_info_id"`
	ToUserID    int64  `json:"video_id"`
	Content     string `json:"content"`
	PublishDate int64  `json:"publish_date"`
}

func (Message) TableName() string {
	return "messages"
}
func CreateMessage(ctx context.Context, message *Message) error {
	return DB.WithContext(ctx).Create(message).Error
}

func QuerryMessageByID(ctx context.Context, userID, toUserID, ttl int64) (messages []*Message, err error) {
	messages = make([]*Message, 0)
	err = DB.WithContext(ctx).Where("publish_date > ?", ttl).
		Where("user_id = ? and to_user_id = ?", userID, toUserID).
		Find(&messages).Error
	return messages, err
}

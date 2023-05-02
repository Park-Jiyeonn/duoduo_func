package db

import (
	"context"
	"simple_tiktok/dal/db/model"
)

func CreateMessage(ctx context.Context, message *model.Message) error {
	return DB.WithContext(ctx).Create(message).Error
}

func QuerryMessageByID(ctx context.Context, userID, toUserID, ttl int64) (messages []*model.Message, err error) {
	messages = make([]*model.Message, 0)
	err = DB.WithContext(ctx).Where("publish_date > ?", ttl).
		Where("user_id = ? and to_user_id = ?", userID, toUserID).
		Find(&messages).Error
	return messages, err
}

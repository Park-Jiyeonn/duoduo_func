package db

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"simple_tiktok/dal/db/model"
)

func CreateMessage(ctx context.Context, message *model.Message) error {
	return DB.WithContext(ctx).Create(message).Error
}

func QuerryMessageByID(ctx context.Context, userID, toUserID, ttl int64) (messages []model.Message, err error) {
	messages = make([]model.Message, 0)
	err = DB.WithContext(ctx).Where("publish_date > ?", ttl).
		Where("user_id = ? and to_user_id = ?", userID, toUserID).
		Find(&messages).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return messages, nil
	}
	return messages, err
}

func GetLastMessageByUid(ctx context.Context, uida, uidb int64) (msg model.Message, err error) {
	tx := DB.WithContext(ctx).Model(model.Message{})
	tx.Where("user_id = ? AND to_user_uid = ?", uida, uidb)
	tx.Or("user_id = ? AND to_user_uid = ?", uidb, uida)
	if err = tx.Order("created_at desc").First(&msg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return msg, nil
		} else {
			return msg, err
		}
	}
	return msg, nil
}

package dal

import "simple_tiktok/pojo"

func CreateMessage(message *pojo.Message) error {
	return DB.Create(message).Error
}

func QuerryMessageByID(userID, toUserID, ttl int64) (messages []pojo.Message, err error) {
	messages = make([]pojo.Message, 0)
	err = DB.Where("publish_date > ?", ttl).
		Where("user_id = ? and to_user_id = ?", userID, toUserID).
		Find(&messages).Error
	return messages, err
}

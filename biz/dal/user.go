package dal

import (
	"simple_tiktok/pojo"
)

func CreateUsers(user *pojo.User) error {
	return DB.Create(user).Error
}

func FindUserByName(userName string) ([]*pojo.User, error) {
	res := make([]*pojo.User, 0)
	err := DB.Where("user_name=?", userName).Find(&res).Error
	return res, err
}
func FindUserByID(userID int64) ([]*pojo.User, error) {
	res := make([]*pojo.User, 0)
	err := DB.Where("id=?", userID).Find(&res).Error
	return res, err
}

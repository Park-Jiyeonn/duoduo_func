package db

import (
	"context"
	"duoduo_fun/dal/db/model"
	"errors"
	"gorm.io/gorm"
	"log"
)

func CreateUser(username, password string) (int, error) {
	// 其他的那些count会自动初始化为0，不用在这里指定了
	user := &model.User{
		Name:     username,
		Password: password,
	}
	err := DB.Model(model.User{}).Create(user).Error
	if err != nil {
		log.Printf("Create User error: %v, user: %+v", err, user)
		return 0, err
	}
	return int(user.ID), nil
}

func GetUserById(ctx context.Context, userId int64) (*model.User, error) {
	var user model.User
	err := DB.WithContext(ctx).Model(model.User{}).Where("id = ?", userId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("Get User by Id error: %v", err)
		return &user, err
	}
	return &user, nil
}

func GetUserByName(name string) (*model.User, error) {
	var user model.User
	err := DB.Model(model.User{}).Where("name = ?", name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return &user, err
	}
	return &user, nil
}

func UpdateUser(uid int64, userMap *map[string]interface{}) (err error) {
	if err = DB.Model(model.User{}).Where("id = ?", uid).Updates(&userMap).Error; err != nil {
		return err
	}
	return nil
}

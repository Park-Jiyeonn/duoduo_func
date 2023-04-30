package db

import (
	"context"
	"log"
	"simple_tiktok/dal/db/model"
)

func CreateUser(ctx context.Context, username, password string) (int64, error) {
	// 其他的那些count会自动初始化为0，不用在这里指定了
	user := &model.User{
		Name:     username,
		Password: password,
	}
	err := DB.WithContext(ctx).Model(model.User{}).Create(user).Error
	if err != nil {
		log.Printf("Create User error: %v, user: %+v", err, user)
		return 0, err
	}
	return int64(user.ID), nil
}

func GetUserById(ctx context.Context, userId int64) (*model.User, error) {
	var user model.User
	err := DB.WithContext(ctx).Model(model.User{}).Where("id = ?", userId).First(&user).Error
	if err != nil {
		log.Printf("Get User by Id error: %v", err)
		return nil, err
	}
	return &user, nil
}

func GetUserByName(ctx context.Context, name string) (*model.User, error) {
	var user model.User
	err := DB.WithContext(ctx).Model(model.User{}).Where("name = ?", name).First(&user).Error
	if err != nil {
		log.Printf("Get User by Id error: %v", err)
		return nil, err
	}
	return &user, nil
}

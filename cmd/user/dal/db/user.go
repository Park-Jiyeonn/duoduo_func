package db

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"simple_tiktok/util/consts"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) TableName() string {
	return consts.UserTableName
}

func GetUser(ctx context.Context, userID int64) (*User, error) {
	user := &User{}
	err := DB.WithContext(ctx).Where("id = ?", userID).First(user).Error;
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with ID %d", userID)
		}
		return nil, fmt.Errorf("failed to get user with ID %d: %w", userID, err)
	}
	return user, nil
}

// CreateUser create user info
func CreateUser(ctx context.Context, user *User) (userID int64, err error) {
	err = DB.WithContext(ctx).Create(user).Error
	return int64(user.ID), err
}

// QueryUser query list of user info
func QueryUser(ctx context.Context, userName string) ([]*User, error) {
	res := make([]*User, 0)
	if err := DB.WithContext(ctx).Where("username = ?", userName).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

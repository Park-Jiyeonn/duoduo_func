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

type Video struct {
	gorm.Model
	UserName  string `json:"username"` // 上传者的用户信息id
	Title     string `json:"title"`    // 视频标题
	VideoPath string `json:"play_url"` // 视频存储位置
	CoverPath string `json:"cover_url"`
}

func (v *Video) TableName() string {
	return consts.VideoTableName
}

func (u *User) TableName() string {
	return consts.UserTableName
}

func GetUser(ctx context.Context, userID int64) (*User, error) {
	user := &User{}
	err := DB.WithContext(ctx).Where("id = ?", userID).First(user).Error
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

// QueryUserByName query list of user info
func QueryUserByName(ctx context.Context, userName string) ([]*User, error) {
	res := make([]*User, 0)
	if err := DB.WithContext(ctx).Where("username = ?", userName).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func QueryUserById(ctx context.Context, userId int64) ([]*User, error) {
	res := make([]*User, 0)
	if err := DB.WithContext(ctx).Where("id = ?", userId).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func CreateVedio(ctx context.Context, vedio *Video) error {
	return DB.WithContext(ctx).Create(vedio).Error
}

func FindVideoByUserID(ctx context.Context, userName string) ([]*Video, error) {
	res := make([]*Video, 0)
	err := DB.WithContext(ctx).Where("user_name=?", userName).Find(&res).Error
	return res, err
}
func FindVideoAll(ctx context.Context) ([]*Video, error) {
	res := make([]*Video, 0)
	err := DB.WithContext(ctx).Order("created_at desc").Find(&res).Error
	return res, err
}

func FindVideoByVideoID(ctx context.Context, VideoID int64) (*Video, error) {
	var res *Video
	err := DB.WithContext(ctx).Where("id=?", VideoID).Find(&res).Error
	return res, err
}

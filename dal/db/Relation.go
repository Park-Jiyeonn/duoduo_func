package db

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"simple_tiktok/dal/db/model"
)

func GetFollowList(ctx context.Context, uid int64) ([]int64, error) {
	var userIds []int64
	err := DB.WithContext(ctx).Model(model.Relation{}).
		Select("to_user_id").
		Where("user_id = ? AND action = 1", uid).
		Find(&userIds).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userIds, nil
		}
		return nil, err
	}
	return userIds, nil
}

func GetFollowerList(ctx context.Context, uid int64) ([]int64, error) {
	var userIds []int64
	err := DB.WithContext(ctx).Model(model.Relation{}).
		Select("user_id").
		Where("to_user_id = ? AND action = 1", uid).
		Find(&userIds).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userIds, nil
		}
		return nil, err
	}
	return userIds, nil
}

func IsFollowed(ctx context.Context, uid, toUserId int64) (bool, error) {
	var ret bool
	err := DB.WithContext(ctx).Model(model.Relation{}).
		Select("action").
		Where("user_id = ? AND to_user_id = ?", uid, toUserId).
		First(&ret).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return ret, nil
}

func UpdateAndInsertRelation(ctx context.Context, uid, toUserId int64, action bool) error {
	relation := &model.Relation{
		UserId:   uid,
		ToUserId: toUserId,
		Action:   action,
	}
	err := DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "to_user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"action"}),
		}).Create(&relation).Error
	if err != nil {
		return err
	}
	return nil
}

package db

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"simple_tiktok/dal/db/model"
)

func GetUserLikeRecords(ctx context.Context, uid int64) ([]int64, error) {
	var likeVideoIds []int64
	err := DB.WithContext(ctx).Model(model.Like{}).
		Select("video_id").
		Where("user_id = ? AND action=1", uid).
		Find(&likeVideoIds).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return likeVideoIds, nil
		}
		return nil, err
	}
	return likeVideoIds, nil
}

func HasLiked(ctx context.Context, uid, vid int64) (bool, error) {
	var ret bool
	err := DB.WithContext(ctx).Model(model.Like{}).
		Select("action").
		Where("user_id=? AND video_id=?", uid, vid).
		First(&ret).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return ret, nil
}

func UpdateAndInsertLikeRecord(ctx context.Context, uid, vid int64, action bool) error {
	record := &model.Like{
		UserId:  uid,
		VideoId: vid,
		Action:  action,
	}
	err := DB.WithContext(ctx).Model(model.Like{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "video_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"action"}),
		}).Create(&record).Error
	if err != nil {
		return err
	}
	return nil
}

package db

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"simple_tiktok/dal/db/model"
)

func GetCommentList(ctx context.Context, vid int64) ([]model.Comment, error) {
	res := make([]model.Comment, 0)
	err := DB.WithContext(ctx).
		Model(model.Comment{}).
		Where("vid = ?", vid).Find(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res, nil
		}
		return nil, err
	}
	return res, nil
}

func CreateComment(ctx context.Context, comment *model.Comment) error {
	return DB.WithContext(ctx).Model(model.Comment{}).
		Create(comment).Error
}

func DeleteCommentByID(ctx context.Context, id int64) error {
	return DB.WithContext(ctx).Model(&model.Comment{}).
		Where("id = ?", id).
		Delete(nil).Error
}

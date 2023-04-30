package db

import (
	"context"
	"gorm.io/gorm"
	"simple_tiktok/dal/db/model"
)

func CreateVideo(ctx context.Context, playUrl string, coverUrl string, title string, uid int64) (err error) {
	video := &model.Video{
		Title:    title,
		PlayUrl:  playUrl,
		CoverUrl: coverUrl,
		UserId:   uid,
	}
	tx := DB.Begin().WithContext(ctx)
	if err = tx.Create(video).Error; err != nil {
		tx.Rollback()
	}
	if err = tx.Model(model.User{}).Where("id = ?", uid).UpdateColumn("work_count", gorm.Expr("work_count+1")).Error; err != nil {
		tx.Rollback()
	}
	tx.Commit()
	return
}

func GetVideoByUserId(ctx context.Context, uid int64) ([]*model.Video, error) {
	var videos []*model.Video
	err := DB.WithContext(ctx).Model(model.Video{}).Where("user_id = ?", uid).Find(videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func MGetByTime(ctx context.Context, latestTime int64) ([]*model.Video, error) {
	tx := DB.WithContext(ctx).Model(model.Video{}).Where("unix_timestamp(created_at) <= ?", latestTime)
	var videos []*model.Video
	err := tx.Order("created_at desc").Limit(30).Find(videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func GetVideoByVideoId(ctx context.Context, vid int64) (*model.Video, error) {
	var video *model.Video
	err := DB.WithContext(ctx).Model(model.Video{}).Where("id = ?", vid).First(video).Error
	if err != nil {
		return nil, err
	}
	return video, nil
}

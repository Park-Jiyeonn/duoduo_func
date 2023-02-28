package db

import (
	"context"
	"gorm.io/gorm"
	"simple_tiktok/cmd/base/dal/db"
)

type Comment struct {
	gorm.Model
	UserID      int64  `json:"user_id"`
	VideoID     int64  `json:"video_id"`
	Content     string `json:"content"`
	PublishDate string `json:"publish_date"`
}

func (Comment) TableName() string {
	return "comments"
}

func FindVideoByVideoID(ctx context.Context, VideoID int64) (*db.Video, error) {
	var res *db.Video
	err := DB.WithContext(ctx).Where("id=?", VideoID).Find(&res).Error
	return res, err
}
func QueryUserByName(ctx context.Context, userName string) ([]*db.User, error) {
	res := make([]*db.User, 0)
	if err := DB.WithContext(ctx).Where("username = ?", userName).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// CreateComment 添加评论
func CreateComment(ctx context.Context, comment *Comment) error {
	return DB.WithContext(ctx).Create(comment).Error
}

// DeleteCommentByID 删除评论
func DeleteCommentByID(ctx context.Context, id int64) error {
	return DB.WithContext(ctx).Where("id = ?", id).Delete(&Comment{}).Error
}

// QueryCommentByVideoID 根据视频ID查询该视频的所有评论
func QueryCommentByVideoID(ctx context.Context, videoId int64) (comments []*Comment, err error) {
	comments = make([]*Comment, 0)
	err = DB.WithContext(ctx).Where("video_id = ?", videoId).Find(&comments).Error
	return comments, err
}

// QueryCommentsCount 查询有多少条评论
func QueryCommentsCount(ctx context.Context, videoId int64) (int64, error) {
	var count int64
	result := DB.WithContext(ctx).Model(&Comment{}).Where("video_id = ?", videoId).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func QueryUserById(ctx context.Context, userId int64) ([]*db.User, error) {
	res := make([]*db.User, 0)
	if err := DB.WithContext(ctx).Where("id = ?", userId).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

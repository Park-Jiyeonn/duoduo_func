package dal

import (
	"simple_tiktok/pojo"
)

// CreateComment 添加评论
func CreateComment(comment *pojo.Comment) error {
	return DB.Create(comment).Error
}

// DeleteCommentByID 删除评论
func DeleteCommentByID(id int64) error {
	return DB.Where("id = ?", id).Delete(&pojo.Comment{}).Error
}

// QueryCommentByVideoID 根据视频ID查询该视频的所有评论
func QueryCommentByVideoID(videoId int64) (comments []pojo.Comment, err error) {
	comments = make([]pojo.Comment, 0)
	err = DB.Where("video_id = ?", videoId).Find(&comments).Error
	return comments, err
}

// 查询有多少条评论
func QueryCommentsCount(videoId int64) (int64, error) {
	var count int64
	result := DB.Model(&pojo.Comment{}).Where("video_id = ?", videoId).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

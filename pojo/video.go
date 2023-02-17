package pojo

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	UserInfoID int64  `json:"user_info_id"` // 上传者的用户信息id
	Title      string `json:"title"`        // 视频标题
	VideoPath  string `json:"play_url"`     // 视频存储位置
	CoverPath  string `json:"cover_url"`    // 封面存储位置
}

func (Video) TableName() string {
	return "videos"
}

package pojo

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	UserName  string `json:"user_name"` // 上传者的用户信息id
	Title     string `json:"title"`     // 视频标题
	VideoPath string `json:"play_url"`  // 视频存储位置
	CoverPath string `json:"cover_url"`
}

func (Video) TableName() string {
	return "videos"
}

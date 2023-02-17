package dal

import "simple_tiktok/pojo"

func CreateVedio(vedio *pojo.Video) error {
	return DB.Create(vedio).Error
}

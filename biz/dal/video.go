package dal

import "simple_tiktok/pojo"

func CreateVedio(vedio *pojo.Video) error {
	return DB.Create(vedio).Error
}

func FindVideoByName(userName string) ([]*pojo.Video, error) {
	res := make([]*pojo.Video, 0)
	err := DB.Where("user_name=?", userName).Find(&res).Error
	return res, err
}
func FindVideoAll() ([]*pojo.Video, error) {
	res := make([]*pojo.Video, 0)
	err := DB.Find(&res).Error
	return res, err
}

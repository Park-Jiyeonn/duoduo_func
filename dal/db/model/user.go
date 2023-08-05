package model

import (
	"gorm.io/gorm"
)

// User Gorm Data structures
type User struct {
	gorm.Model
	Name          string `gorm:"column:name;index:uni_name,unique;type:varchar(32);not null"  json:"name" redis:"name"`
	Password      string `gorm:"column:password;type:varchar(255);not null"  json:"password" redis:"password"`
	WorkCount     int    `gorm:"column:work_count;default:0"  json:"work_count" redis:"work_count"`
	FavoriteCount int    `gorm:"column:favorite_count;default:0" json:"favorite_count" redis:"favorite_count"`
}

func (User) TableName() string {
	return "user"
}

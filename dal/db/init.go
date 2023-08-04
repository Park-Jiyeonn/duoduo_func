package db

import (
	"duoduo_fun/dal/db/model"
	"duoduo_fun/pkg/consts"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	if DB != nil {
		return
	}
	var err error
	DB, err = gorm.Open(mysql.Open(consts.MySQLDefaultDSN),
		&gorm.Config{
			PrepareStmt: true,
		},
	)
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(
		&model.User{},
		&model.Video{},
		&model.Comment{},
		&model.Like{},
	)
	if err != nil {
		panic(err)
	}
}

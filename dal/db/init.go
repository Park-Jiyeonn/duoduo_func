package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/logging/logrus"
	"gorm.io/plugin/opentelemetry/tracing"
	"simple_tiktok/dal/db/model"
	"simple_tiktok/pkg/consts"
	"time"
)

var DB *gorm.DB

func Init() {
	var err error
	gormlogrus := logger.New(
		logrus.NewWriter(),
		logger.Config{
			SlowThreshold: time.Millisecond,
			Colorful:      false,
			LogLevel:      logger.Info,
		},
	)
	//开启 PrepareStmt 功能，可以加快查询速度；
	//配置 GORM 的 Logger 为 gormlogrus，也就是使用 logrus 包输出 SQL 执行日志；
	// 这里连接的端口号，是docker中的端口号！
	DB, err = gorm.Open(mysql.Open(consts.MySQLDefaultDSN),
		&gorm.Config{
			PrepareStmt: true,
			Logger:      gormlogrus,
		},
	)
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&model.User{}, &model.Video{}, &model.Message{}, &model.Comment{})
	if err != nil {
		panic(err)
	}

	//if !DB.Migrator().HasIndex(&User{}, "idx_user_name") {
	//	err = DB.Migrator().CreateIndex(&User{}, "idx_user_name")
	//	if err != nil {
	//		panic(err)
	//	}
	//}
	//
	//if !DB.Migrator().HasIndex(&Video{}, "idx_video_username") {
	//	err = DB.Migrator().CreateIndex(&Video{}, "idx_video_username")
	//	if err != nil {
	//		panic(err)
	//	}
	//}
	//
	//if !DB.Migrator().HasIndex(&Comment{}, "idx_comment_user_id") {
	//	err = DB.Migrator().CreateIndex(&Comment{}, "idx_comment_user_id")
	//	if err != nil {
	//		panic(err)
	//	}
	//}
	//
	//if !DB.Migrator().HasIndex(&Message{}, "idx_message_user_id") {
	//	err = DB.Migrator().CreateIndex(&Message{}, "idx_message_user_id")
	//	if err != nil {
	//		panic(err)
	//	}
	//}

	//将 tracing.NewPlugin() 注册到 DB 实例中，用于开启 GORM 链路追踪。
	if err := DB.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}
}

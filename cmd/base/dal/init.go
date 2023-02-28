package dal

import (
	"simple_tiktok/cmd/base/dal/db"
	"simple_tiktok/cmd/base/dal/redis"
)

func Init() {
	db.Init()
	redis.InitRedis()
}

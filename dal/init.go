package dal

import (
	"simple_tiktok/dal/db"
	"simple_tiktok/dal/redis"
)

func Init() {
	db.Init()
	redis.InitRedis()
}

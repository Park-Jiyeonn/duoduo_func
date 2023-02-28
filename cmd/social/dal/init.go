package dal

import (
	"simple_tiktok/cmd/social/dal/db"
	"simple_tiktok/cmd/social/dal/redis"
)

func Init() {
	db.Init()
	redis.InitRedis()
}

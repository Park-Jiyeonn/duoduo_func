package dal

import (
	"simple_tiktok/cmd/interact/dal/db"
	"simple_tiktok/cmd/interact/dal/redis"
)

func Init() {
	db.Init()
	redis.InitRedis()
}

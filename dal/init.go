package dal

import (
	"duoduo_fun/dal/db"
	"duoduo_fun/dal/redis"
)

func Init() {
	db.Init()
	redis.InitRedis()
}

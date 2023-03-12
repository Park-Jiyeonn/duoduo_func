package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var Rs *redis.Client

func InitRedis() {
	// 这里连接的端口号，是docker中的端口号！
	Rs = redis.NewClient(&redis.Options{
		Addr:     "localhost:5070",
		Password: "",
		DB:       0,
	})

	_, err := Rs.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}

package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var Rs *redis.Client

func InitRedis() {
	Rs = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := Rs.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}

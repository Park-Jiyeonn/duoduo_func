package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var Rs *redis.Client

func InitRedis() {
	Rs = redis.NewClient(&redis.Options{
		Addr:     "localhost:5070",
		Password: "",
		DB:       0,
	})

	_, err := Rs.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	SyncDataToDB()
}

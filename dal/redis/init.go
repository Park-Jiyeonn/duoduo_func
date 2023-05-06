package redis

import (
	"context"
	"github.com/go-co-op/gocron"
	"github.com/redis/go-redis/v9"
	"time"
)

var Rs *redis.Client

func InitRedis() {
	if Rs != nil {
		return
	}
	Rs = redis.NewClient(&redis.Options{
		Addr:     "localhost:5070",
		Password: "",
		DB:       0,
	})

	_, err := Rs.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	s := gocron.NewScheduler(time.Local)
	_, _ = s.Every(5).Minutes().Do(SyncDataToDB)
	//s.StartBlocking()
}

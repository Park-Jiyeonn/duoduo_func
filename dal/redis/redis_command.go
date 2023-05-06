package redis

import (
	"context"
	"encoding/json"
	"time"
)

// ========通用操作==========

func Exists(ctx context.Context, key ...string) int64 {
	return Rs.Exists(ctx, key...).Val()
}

// =======String操作==============

func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (ok bool) {
	if value == nil {
		return false
	}
	bs, err := json.Marshal(value)
	if err != nil {
		//Log.Warnf("marshal value err: %v", err)
		return false
	}
	if expiration == 0 {
		expiration = time.Hour
	}
	if status := Rs.Set(ctx, key, bs, expiration); status.Err() != nil {
		//Log.Errorf("set %v to redis err: %v", key, status.Err())
		return false
	}
	return true
}

func Incr(ctx context.Context, key string) (cnt int64, err error) {
	cnt, err = Rs.Incr(ctx, key).Result()
	if err != nil {
		//Log.Warnf("Incr %v redis err: %v", key, err)
		return
	}
	return
}

// ==========Hash操作============

func HSet(ctx context.Context, key string, value interface{}) (ok bool) {
	if status := Rs.HSet(ctx, key, value); status.Err() != nil {
		return false
	}
	return true
}

// HGetAll 返回值为map，需要调用方自行转换
func HGetAll(ctx context.Context, key string) map[string]string {
	return Rs.HGetAll(ctx, key).Val()
}

func HMGet(ctx context.Context, key string, field ...string) []interface{} {
	return Rs.HMGet(ctx, key, field...).Val()
}

func HIncr(ctx context.Context, key, field string, incr int64) (ok bool) {
	if status := Rs.HIncrBy(ctx, key, field, incr); status.Err() != nil {
		return false
	}
	return true
}

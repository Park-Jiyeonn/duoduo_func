package redis

import "context"

func Exists(ctx context.Context, key ...string) int64 {
	return Rs.Exists(ctx, key...).Val()
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

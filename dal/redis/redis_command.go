package redis

import "context"

func Exists(ctx context.Context, key ...string) int64 {
	return Rs.Exists(ctx, key...).Val()
}

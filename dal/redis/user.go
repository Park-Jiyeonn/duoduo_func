package redis

import (
	"context"
	"fmt"
	"simple_tiktok/dal/db/model"
)

func IsExists(ctx context.Context, uids ...int64) int64 {
	keys := make([]string, len(uids))
	for i, val := range uids {
		keys[i] = fmt.Sprintf("user_info_%d", val)
	}
	return Exists(ctx, keys...)
}

func GetUserInfo(ctx context.Context, uid int64) (*model.User, error) {
	
}

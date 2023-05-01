package redis

import (
	"context"
	"simple_tiktok/dal/db/model"
	"simple_tiktok/pkg/consts"
	"strconv"
)

func UserIsExists(ctx context.Context, uids ...int64) int64 {
	keys := make([]string, len(uids))
	for i, val := range uids {
		keys[i] = consts.GetUserInfoKey(val)
	}
	return Exists(ctx, keys...)
}

func GetNum(f map[string]string, s string) int64 {
	ret, _ := strconv.Atoi(f[s])
	return int64(ret)
}

func GetUserInfo(ctx context.Context, uid int64) (*model.User, error) {
	key := consts.GetUserInfoKey(uid)
	value := HGetAll(ctx, key)
	var user *model.User
	user.Name = value["name"]
	user.ID = uint(GetNum(value, "id"))
	user.FollowCount = GetNum(value, "follow_count")
	user.FollowerCount = GetNum(value, "follower_count")
	user.FavoriteCount = GetNum(value, "favorite_count")
	user.WorkCount = GetNum(value, "work_count")
	return user, nil
}

func SetUserInfo(ctx context.Context, user *model.User) bool {
	key := consts.GetUserInfoKey(int64(user.ID))
	return HSet(ctx, key, user)
}

// 1. 根据key和字段名查找值，key或field不存在时，对应的值返回nil，需要调用方自己判断
// 2. 返回的类型都为string，调用方自行转换
func GetUserFields(ctx context.Context, userID int64, field ...string) []interface{} {
	return HMGet(ctx, consts.GetUserInfoKey(userID), field...)
}

func IncrUserField(ctx context.Context, userID int64, field string, incr int64) (ok bool) {
	return HIncr(ctx, consts.GetUserInfoKey(userID), field, incr)
}

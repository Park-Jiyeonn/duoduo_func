package redis

import (
	"context"
	"duoduo_fun/pkg/consts"
	"strconv"
)

func LikeIsExists(ctx context.Context, uids ...int) int64 {
	keys := make([]string, len(uids))
	for i, uid := range uids {
		keys[i] = consts.GetUserLikeListKey(uid)
	}
	return Exists(ctx, keys...)
}

func IsLike(ctx context.Context, uid, vid int) bool {
	res := HMGet(ctx, consts.GetUserLikeListKey(uid), strconv.FormatInt(int64(vid), 10))
	if res[0] == nil || res[0].(string) == "0" {
		return false
	}
	return true
}

func SetFavoriteList(ctx context.Context, userID int, kv ...string) bool {
	return HSet(ctx, consts.GetUserLikeListKey(userID), kv)
}

func FavoriteAction(ctx context.Context, uid, vid int, action int) bool {
	return HIncr(ctx, consts.GetUserLikeListKey(uid), strconv.FormatInt(int64(vid), 10), int64(action))
}

func GetAllUserLikes(ctx context.Context, uid int) (userLikes []int) {
	userLikes = make([]int, 0)
	res := HGetAll(ctx, consts.GetUserLikeListKey(uid))
	for k, v := range res {
		vid, _ := strconv.ParseInt(k, 10, 64)
		action, _ := strconv.ParseInt(v, 10, 64)
		if action == 1 {
			userLikes = append(userLikes, int(vid))
		}
	}
	return
}

func GetFavoriteList(ctx context.Context, userID int) map[string]string {
	return HGetAll(ctx, consts.GetUserLikeListKey(userID))
}

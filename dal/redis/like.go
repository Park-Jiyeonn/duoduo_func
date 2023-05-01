package redis

import (
	"context"
	"simple_tiktok/pkg/consts"
	"strconv"
)

func LikeIsExists(ctx context.Context, uids ...int64) int64 {
	keys := make([]string, len(uids))
	for i, uid := range uids {
		keys[i] = consts.GetUserLikeKey(uid)
	}
	return Exists(ctx, keys...)
}

func IsLike(ctx context.Context, uid, vid int64) bool {
	res := HMGet(ctx, consts.GetUserLikeKey(uid), strconv.FormatInt(vid, 10))
	if res[0] == nil || res[0].(string) == "0" {
		return false
	}
	return true
}

func SetFavoriteList(ctx context.Context, userID int64, kv ...string) bool {
	return HSet(ctx, consts.GetUserLikeKey(userID), kv)
}

func FavoriteAction(ctx context.Context, uid, vid int64, action int64) bool {
	return HIncr(ctx, consts.GetUserLikeKey(uid), strconv.FormatInt(vid, 10), action)
}

func GetAllUserLikes(ctx context.Context, uid int64) (userLikes []*int64) {
	userLikes = make([]*int64, 0)
	res := HGetAll(ctx, consts.GetUserLikeKey(uid))
	for k, v := range res {
		vid, _ := strconv.ParseInt(k, 10, 64)
		action, _ := strconv.ParseInt(v, 10, 64)
		if action == 1 {
			userLikes = append(userLikes, &vid)
		}
	}
	return
}

func GetFavoriteList(ctx context.Context, userID int64) map[string]string {
	return HGetAll(ctx, consts.GetUserLikeKey(userID))
}

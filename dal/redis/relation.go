package redis

import (
	"context"
	"simple_tiktok/pkg/consts"
	"strconv"
)

func FollowIsExists(ctx context.Context, uids ...int64) int64 {
	keys := make([]string, len(uids))
	for i, val := range uids {
		keys[i] = consts.GetUserFollowListKey(val)
	}
	return Exists(ctx, keys...)
}

func FollowerIsExists(ctx context.Context, uids ...int64) int64 {
	keys := make([]string, len(uids))
	for i, uid := range uids {
		keys[i] = consts.GetUserFollowerListKey(uid)
	}
	return Exists(ctx, keys...)
}

func IsFollow(ctx context.Context, userId, toUserId int64) bool {
	res := HMGet(ctx, consts.GetUserFollowListKey(userId), strconv.FormatInt(toUserId, 10))
	if res[0] == nil || res[0].(string) == "0" {
		return false
	}
	return true
}

func SetFollowList(ctx context.Context, userID int64, kv ...string) bool {
	return HSet(ctx, consts.GetUserFollowListKey(userID), kv)
}

func SetFollowerList(ctx context.Context, userID int64, kv ...string) bool {
	return HSet(ctx, consts.GetUserFollowerListKey(userID), kv)
}

func FollowAction(ctx context.Context, userId, toUserId int64, action int64) bool {
	b1 := HIncr(ctx, consts.GetUserFollowListKey(userId), strconv.FormatInt(toUserId, 10), action)
	b2 := HIncr(ctx, consts.GetUserFollowerListKey(toUserId), strconv.FormatInt(userId, 10), action)
	return b1 && b2
}

func GetFollowList(ctx context.Context, userID int64) (followList []int64) {
	followList = make([]int64, 0)
	res := HGetAll(ctx, consts.GetUserFollowListKey(userID))
	for k, v := range res {
		uid, _ := strconv.ParseInt(k, 10, 64)
		action, _ := strconv.ParseInt(v, 10, 64)
		if action == 1 {
			followList = append(followList, uid)
		}
	}
	return
}

func GetFollowerList(ctx context.Context, userID int64) (followerList []int64) {
	followerList = make([]int64, 0)
	res := HGetAll(ctx, consts.GetUserFollowerListKey(userID))
	for k, v := range res {
		uid, _ := strconv.ParseInt(k, 10, 64)
		action, _ := strconv.ParseInt(v, 10, 64)
		if action == 1 {
			followerList = append(followerList, uid)
		}
	}
	return
}

func GetFollowFullList(ctx context.Context, userID int64) map[string]string {
	return HGetAll(ctx, consts.GetUserFollowListKey(userID))
}

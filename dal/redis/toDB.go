package redis

import (
	"context"
	"simple_tiktok/dal/db"
	"simple_tiktok/pkg/consts"
	"strconv"
)

func SyncDataToDB() {
	ctx := context.Background()
	go UpdateUserMsgToDB(ctx)
	go UpdateUserLikeListToDB(ctx)
	go UpdateUserFollowListToDB(ctx)
	go UpdateVideoMsgToDB(ctx)
}

func UpdateUserMsgToDB(ctx context.Context) {
	var cursor uint64
	keys, _, err := Rs.Scan(ctx, cursor, "user_info_*", 0).Result()
	if err != nil {
		return
	}
	for _, key := range keys {
		uid := consts.GetIDFromUserMsgKey(key)
		user, err := GetUserInfo(ctx, uid)
		if err != nil {
			continue
		}
		userMap := map[string]interface{}{
			"id":             user.ID,
			"name":           user.Name,
			"follow_count":   user.FollowCount,
			"follower_count": user.FollowerCount,
			"work_count":     user.WorkCount,
			"favorite_count": user.FavoriteCount,
		}
		err = db.UpdateUser(ctx, uid, &userMap)
		if err != nil {
			continue
		}
	}
	return
}

func UpdateUserLikeListToDB(ctx context.Context) {
	var cursor uint64
	keys, _, err := Rs.Scan(ctx, cursor, "user_like_list_*", 0).Result()
	if err != nil {
		return
	}
	for _, key := range keys {
		uid := consts.GetIDFromUserLikeListKey(key)
		list := GetFavoriteList(ctx, uid)
		for k, v := range list {
			vid, err := strconv.ParseInt(k, 10, 64) //nolint: staticcheck
			action, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				continue
			}
			err = db.UpdateAndInsertLikeRecord(ctx,
				uid,
				vid,
				action == 1,
			)
			if err != nil {
				continue
			}
		}
	}
	return
}

func UpdateUserFollowListToDB(ctx context.Context) {
	var cursor uint64
	keys, _, err := Rs.Scan(ctx, cursor, "user_follow_list*", 0).Result()
	if err != nil {
		return
	}
	for _, key := range keys {
		uid := consts.GetIDFromUserFollowListKey(key)
		list := GetFollowFullList(ctx, uid)
		for k, v := range list {
			uid2, err := strconv.ParseInt(k, 10, 64) //nolint: staticcheck
			action, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				continue
			}
			err = db.UpdateAndInsertRelation(ctx,
				uid,
				uid2,
				action == 1,
			)
			if err != nil {
				continue
			}
		}
	}
	return
}

func UpdateVideoMsgToDB(ctx context.Context) {
	var cursor uint64
	keys, _, err := Rs.Scan(ctx, cursor, "video_message_*", 0).Result()
	if err != nil {
		return
	}
	for _, key := range keys {
		vid := consts.GetIDFromVideoMsgKey(key)
		video, err := GetVideoMessage(ctx, vid)
		if err != nil {
			continue
		}
		videoMap := map[string]interface{}{
			"id":             video.ID,
			"play_url":       video.PlayUrl,
			"cover_url":      video.CoverUrl,
			"favorite_count": video.FavoriteCount,
			"comment_count":  video.CommentCount,
			"title":          video.Title,
			"user_id":        video.UserId,
		}
		err = db.UpdateVideo(ctx, vid, &videoMap)
		if err != nil {
			continue
		}
	}
	return
}

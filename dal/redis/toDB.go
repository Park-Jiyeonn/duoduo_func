package redis

import (
	"context"
	"duoduo_fun/dal/db"
	"duoduo_fun/pkg/consts"
	"strconv"
)

func SyncDataToDB() {
	ctx := context.Background()
	go UpdateUserMsgToDB(ctx)
	go UpdateUserLikeListToDB(ctx)
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
		user, err := GetUserInfo(ctx, int(uid))
		if err != nil {
			continue
		}
		userMap := map[string]interface{}{
			"id":             user.ID,
			"name":           user.Name,
			"work_count":     user.WorkCount,
			"favorite_count": user.FavoriteCount,
		}
		err = db.UpdateUser(uid, &userMap)
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
		list := GetFavoriteList(ctx, int(uid))
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

func UpdateVideoMsgToDB(ctx context.Context) {
	var cursor uint64
	keys, _, err := Rs.Scan(ctx, cursor, "video_message_*", 0).Result()
	if err != nil {
		return
	}
	for _, key := range keys {
		vid := consts.GetIDFromVideoMsgKey(key)
		video, err := GetVideoMessage(ctx, int(vid))
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

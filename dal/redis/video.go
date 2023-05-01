package redis

import (
	"context"
	"simple_tiktok/dal/db/model"
	"simple_tiktok/pkg/consts"
)

func VideoIsExists(ctx context.Context, vids ...int64) int64 {
	keys := make([]string, len(vids))
	for i, vid := range vids {
		keys[i] = consts.GetVideoMsgKey(vid)
	}
	return Exists(ctx, keys...)
}

func SetVideoMessage(ctx context.Context, video *model.Video) (ok bool) {
	return HSet(ctx, consts.GetVideoMsgKey(int64(video.ID)), video)
}

func IncrVideoField(ctx context.Context, vid int64, field string, incr int64) (ok bool) {
	return HIncr(ctx, consts.GetVideoMsgKey(vid), field, incr)
}

func GetVideoFields(ctx context.Context, vid int64, field ...string) []interface{} {
	return HMGet(ctx, consts.GetVideoMsgKey(vid), field...)
}

func GetVideoMessage(ctx context.Context, vid int64) (video *model.Video, err error) {
	key := consts.GetVideoMsgKey(vid)
	value := HGetAll(ctx, key)
	video.UserId = GetNum(value, "user_id")
	video.Title = value["title"]
	video.FavoriteCount = GetNum(value, "favorite_count")
	video.PlayUrl = value["play_url"]
	video.CommentCount = GetNum(value, "comment_count")
	video.CoverUrl = value["cover_url"]
	return video, nil
}

package redis

import (
	"context"
	"duoduo_fun/dal/db/model"
	"duoduo_fun/pkg/consts"
)

type redisVideo struct {
	Id            int    `json:"id" redis:"id"`
	UserId        int    `json:"user_id" redis:"user_id"`
	PlayUrl       string `json:"play_url" redis:"play_url"`
	CoverUrl      string `json:"cover_url" redis:"cover_url"`
	FavoriteCount int    `json:"favorite_count" redis:"favorite_count"`
	CommentCount  int    `json:"comment_count" redis:"comment_count"`
	Title         string `json:"title" redis:"title"`
}

func VideoIsExists(ctx context.Context, vids ...int) int64 {
	keys := make([]string, len(vids))
	for i, vid := range vids {
		keys[i] = consts.GetVideoMsgKey(vid)
	}
	return Exists(ctx, keys...)
}

func SetVideoMessage(ctx context.Context, video *model.Video) (ok bool) {
	return HSet(ctx, consts.GetVideoMsgKey(int(video.ID)), &redisVideo{
		Id:            int(video.ID),
		UserId:        video.UserId,
		PlayUrl:       video.PlayUrl,
		CoverUrl:      video.CoverUrl,
		FavoriteCount: video.FavoriteCount,
		CommentCount:  video.CommentCount,
		Title:         video.Title,
	})
}

func IncrVideoField(ctx context.Context, vid int, field string, incr int64) (ok bool) {
	return HIncr(ctx, consts.GetVideoMsgKey(vid), field, incr)
}

func GetVideoFields(ctx context.Context, vid int, field ...string) []interface{} {
	return HMGet(ctx, consts.GetVideoMsgKey(vid), field...)
}

func GetVideoMessage(ctx context.Context, vid int) (*model.Video, error) {
	key := consts.GetVideoMsgKey(vid)
	value := HGetAll(ctx, key)
	video := new(model.Video)
	video.UserId = GetNum(value, "user_id")
	video.Title = value["title"]
	video.FavoriteCount = GetNum(value, "favorite_count")
	video.PlayUrl = value["play_url"]
	video.CommentCount = GetNum(value, "comment_count")
	video.CoverUrl = value["cover_url"]
	return video, nil
}

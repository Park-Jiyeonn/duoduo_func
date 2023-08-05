package handler

import (
	"context"
	api_model "duoduo_fun/api/model"
	"duoduo_fun/dal/db"
	"duoduo_fun/dal/db/model"
	"duoduo_fun/dal/redis"
	"duoduo_fun/pkg/errno"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func likeAction(ctx context.Context, req *api_model.LikeRequest) (resp *api_model.LikeResponse, err error) {
	if redis.LikeIsExists(ctx, req.UserId) == 0 {
		likeList, err := db.GetUserLikeRecords(ctx, req.UserId)
		if err != nil {
			resp.StatusCode = 1
			return resp, err
		}
		if len(likeList) > 0 {
			kv := make([]string, 0)
			for _, videoId := range likeList {
				kv = append(kv, strconv.FormatInt(int64(videoId), 10))
				kv = append(kv, "1")
			}
			if !redis.SetFavoriteList(ctx, req.UserId, kv...) {
				resp.StatusCode = 1
				return resp, errno.NewErrNo("Redis设置用户点赞视频缓存出错")
			}
		}
	}
	if redis.UserIsExists(ctx, req.UserId) == 0 {
		user, err := db.GetUserById(ctx, req.UserId)
		if err != nil {
			resp.StatusCode = 1
			return resp, err
		}
		if !redis.SetUserInfo(ctx, user) {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("Redis缓存用户信息出错")
		}
	}
	if redis.VideoIsExists(ctx, req.VideoId) == 0 {
		video, err := db.GetVideoByVideoId(ctx, req.VideoId)
		if err != nil {
			return nil, err
		}
		if !redis.SetVideoMessage(ctx, video) {
			return nil, errno.NewErrNo("Redis缓存视频信息出错")
		}
	}
	res := redis.GetVideoFields(ctx, req.VideoId, "user_id")
	authorID, _ := strconv.ParseInt(res[0].(string), 10, 64)
	if redis.UserIsExists(ctx, int(authorID)) == 0 {
		user, err := db.GetUserById(ctx, int(authorID))
		if err != nil {
			resp.StatusCode = 1
			return resp, err
		}
		if !redis.SetUserInfo(ctx, user) {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("Redis缓存用户信息出错")
		}
	}

	var action int64
	like := redis.IsLike(ctx, req.UserId, req.VideoId)
	if like && req.ActionType == "2" {
		action = -1
	} else if !like && req.ActionType == "1" {
		action = 1
	} else {
		return
	}

	// 更新点赞列表
	if !redis.FavoriteAction(ctx, req.UserId, req.VideoId, int(action)) {
		return resp, errno.NewErrNo("更新点赞列表失败")
	}
	// 更新用户点赞数
	if !redis.IncrUserField(ctx, req.UserId, "favorite_count", action) {
		return resp, errno.NewErrNo("更新用户点赞数")
	}
	// 更新作者获赞数
	//if !redis.IncrUserField(ctx, authorID, "total_favorited", action) {
	//	return resp, errno.NewErrNo("更新点赞列表失败")
	//}
	// 更新视频获赞数
	if !redis.IncrVideoField(ctx, req.VideoId, "favorite_count", action) {
		return resp, errno.NewErrNo("更新视频获赞数")
	}

	resp.StatusCode = 200
	resp.StatusMsg = "success"
	return resp, nil
}
func LikeAction(c *gin.Context) {
	req := &api_model.LikeRequest{}
	err := c.BindJSON(req)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := likeAction(context.Background(), req)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func getLikeList(ctx context.Context, req *api_model.LikeListRequest) (resp *api_model.LikeListResponse, err error) {
	var likeList []int
	if redis.LikeIsExists(ctx, req.UserId) == 0 {
		likeList, err = db.GetUserLikeRecords(ctx, req.UserId)
		if err != nil {
			resp.StatusCode = 1
			return resp, err
		}
		if len(likeList) > 0 {
			kv := make([]string, 0)
			for _, videoId := range likeList {
				kv = append(kv, strconv.FormatInt(int64(videoId), 10))
				kv = append(kv, "1")
			}
			if !redis.SetFavoriteList(ctx, req.UserId, kv...) {
				resp.StatusCode = 1
				return resp, errno.NewErrNo("Redis设置用户点赞视频缓存出错")
			}
		}
	} else {
		likeList = redis.GetAllUserLikes(ctx, req.UserId)
	}

	var videoList []api_model.VideoInfo
	for _, vid := range likeList {
		var video *model.Video
		if redis.VideoIsExists(ctx, vid) == 0 {
			video, err = db.GetVideoByVideoId(ctx, vid)
			if err != nil {
				return nil, err
			}
			redis.SetVideoMessage(ctx, video)
		} else {
			video, err = redis.GetVideoMessage(ctx, vid)
		}

		var user *model.User
		if redis.UserIsExists(ctx, video.UserId) == 0 {
			user, err = db.GetUserById(ctx, video.UserId)
			if err != nil {
				return nil, err
			}
			redis.SetUserInfo(ctx, user)
		} else {
			user, err = redis.GetUserInfo(ctx, video.UserId)
		}

		retVideo := api_model.VideoInfo{
			Id: vid,
			Author: api_model.UserInfo{
				Id:   video.UserId,
				Name: user.Name,
			},
			PlayUrl:       "http://192.168.137.1:8888/data/" + video.PlayUrl,
			CoverUrl:      "http://192.168.137.1:8888/data/" + video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  0,
			IsFavorite:    true,
			Title:         video.Title,
		}
		videoList = append(videoList, retVideo)
	}

	//fmt.Println(videoList)
	resp.VideoList = videoList
	resp.StatusCode = 0
	resp.StatusMsg = "success"
	return resp, nil
}
func GetLikeList(c *gin.Context) {
	req := &api_model.LikeListRequest{}
	err := c.BindJSON(req)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := getLikeList(context.Background(), req)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

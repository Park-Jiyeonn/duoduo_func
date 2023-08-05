package handler

import (
	"context"
	api_model "duoduo_fun/api/model"
	"duoduo_fun/dal/db"
	"duoduo_fun/dal/db/model"
	"duoduo_fun/dal/redis"
	"duoduo_fun/pkg/errno"
	"duoduo_fun/pkg/mq"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func getVideoList(ctx context.Context, req *api_model.FeedRequest) (resp *api_model.FeedResponse, err error) {
	latestTime, err := strconv.Atoi(req.LatestTime)
	if err != nil {
		return resp, errno.NewErrNo("转换数字时间失败")
	}
	videos, err := db.MGetByTime(ctx, int64(latestTime))
	if err != nil {
		resp.StatusCode = 1
		return nil, err
	}

	var videoList []api_model.VideoInfo
	for _, video := range videos {
		//fmt.Println(video.CoverPath)
		var user *model.User
		if redis.UserIsExists(ctx, video.UserId) == 0 {
			user, err = db.GetUserById(ctx, video.UserId)
			if err != nil {
				return nil, err
			}

			redis.SetUserInfo(ctx, user)
		} else {
			user, err = redis.GetUserInfo(ctx, video.UserId)
			if err != nil {
				return nil, err
			}
		}

		if redis.VideoIsExists(ctx, int(video.ID)) == 0 {
			redis.SetVideoMessage(ctx, &video)
		} else {
			latestVideo, err := redis.GetVideoMessage(ctx, int(video.ID))
			if err != nil {
				return nil, errno.NewErrNo("redis获取最新视频信息失败")
			}
			video.FavoriteCount = latestVideo.FavoriteCount
			video.CommentCount = latestVideo.CommentCount
		}

		isLike := false
		if redis.LikeIsExists(ctx, req.UserId) != 0 {
			isLike = redis.IsLike(ctx, req.UserId, int(video.ID))
		} else {
			isLike, err = db.HasLiked(ctx, req.UserId, int(video.ID))
			if err != nil {
				return
			}
			var action int64
			if isLike {
				action = 1
			} else {
				action = 0
			}
			redis.FavoriteAction(ctx, req.UserId, int(video.ID), int(action))
		}

		video := api_model.VideoInfo{
			Id: int(video.ID),
			Author: api_model.UserInfo{
				Id:   int(user.ID),
				Name: user.Name,
			},
			PlayUrl:       "http://192.168.137.1:8888/data/" + video.PlayUrl,
			CoverUrl:      "http://192.168.137.1:8888/data/" + video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    isLike,
			Title:         video.Title,
		}
		videoList = append(videoList, video)
	}

	resp.VideoList = videoList
	resp.StatusMsg = "success"
	return resp, nil
}
func GetVideoList(c *gin.Context) {
	req := &api_model.FeedRequest{}
	err := c.BindJSON(req)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := getVideoList(context.Background(), req)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func publishAction(ctx context.Context, req *api_model.PublishRequest) (resp *api_model.PublishResponse, err error) {
	VideoPath := req.Title + ".mp4"
	CoverPath := req.Title + ".jpg"

	err = mq.Produce(req.UserId, "../../data/"+VideoPath, "../../data/"+CoverPath)
	go mq.Consume(ctx)
	if err != nil {
		resp.StatusCode = 1
		return nil, err
	}
	resp.StatusCode = 2000
	resp.StatusMsg = "success"
	return resp, nil
}
func PublishAction(c *gin.Context) {
	req := &api_model.PublishRequest{}
	err := c.BindJSON(req)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := publishAction(context.Background(), req)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func getPublishList(ctx context.Context, req *api_model.PublishListRequest) (resp *api_model.PublishListResponse, err error) {
	var user = new(model.User)
	if redis.UserIsExists(ctx, req.UserId) == 0 {
		user, err = db.GetUserById(ctx, req.UserId)
		if err != nil {
			return nil, err
		}
		redis.SetUserInfo(ctx, user)
	} else {
		user, err = redis.GetUserInfo(ctx, req.UserId)
		if err != nil {
			return nil, err
		}
	}

	videos, err := db.GetVideosByUserId(ctx, int64(user.ID))

	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("没找到视频相关信息，寄")
	}
	var videoList []api_model.VideoInfo
	for _, video := range videos {
		var isLike bool
		if redis.LikeIsExists(ctx, video.UserId) == 0 {
			userLikes, err := db.GetUserLikeRecords(ctx, req.UserId)
			if err != nil {
				// TODO: 小心缓存击穿
				isLike = false
			}
			// 将查询到数据的加入缓存
			if len(userLikes) > 0 {
				kv := make([]string, 0)
				for _, userLike := range userLikes {
					kv = append(kv, strconv.FormatInt(int64(userLike), 10))
					kv = append(kv, "1")
				}
				if !redis.SetFavoriteList(ctx, req.UserId, kv...) {
					return resp, err
				}
			}
		}
		isLike = redis.IsLike(ctx, req.UserId, int(video.ID))

		video := api_model.VideoInfo{
			Id: int(video.ID),
			Author: api_model.UserInfo{
				Id:   int(user.ID),
				Name: user.Name,
			},
			PlayUrl:       "http://192.168.137.1:8888/data/" + video.PlayUrl,
			CoverUrl:      "http://192.168.137.1:8888/data/" + video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    isLike,
			Title:         video.Title,
		}
		videoList = append(videoList, video)
	}

	resp.VideoList = videoList
	resp.StatusMsg = "success"
	return resp, err
}
func GetPublishList(c *gin.Context) {
	req := &api_model.PublishListRequest{}
	err := c.BindJSON(req)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := getPublishList(context.Background(), req)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

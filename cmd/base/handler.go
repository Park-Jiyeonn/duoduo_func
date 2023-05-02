package main

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"simple_tiktok/dal/db"
	"simple_tiktok/dal/db/model"
	"simple_tiktok/dal/redis"
	base "simple_tiktok/kitex_gen/base"
	"simple_tiktok/pkg/errno"
	util "simple_tiktok/util/ffmpeg"
	"simple_tiktok/util/jwt"
	"strconv"
	"time"
)

// BaseServiceImpl implements the last service interface defined in the IDL.
type BaseServiceImpl struct{}

// UserRegister implements the BaseServiceImpl interface.
func (s *BaseServiceImpl) UserRegister(ctx context.Context, request *base.RegisterRequest) (resp *base.RegisterResponse, err error) {
	// TODO: Your code here...
	resp = new(base.RegisterResponse)

	user, err := db.GetUserByName(ctx, request.Username)
	if user != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("用户已存在")
	}
	EncodePassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("密码加密失败")
	}
	id, err := db.CreateUser(ctx, request.Username, string(EncodePassword))
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("创建用户失败")
	}

	token, err := jwt.GenToken(id, request.Username)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("创建jwt失败")
	}

	resp.StatusCode = 0
	resp.StatusMsg = "success"
	resp.UserId = id
	resp.Token = token
	return resp, nil
}

// UserLogin implements the BaseServiceImpl interface.
func (s *BaseServiceImpl) UserLogin(ctx context.Context, request *base.LoginRequest) (resp *base.LoginResponse, err error) {
	// TODO: Your code here...
	resp = new(base.LoginResponse)
	user, err := db.GetUserByName(ctx, request.Username)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("数据库查询用户存不存在失败！")
	}

	//密码认证
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("用户名或密码错误")
	}

	//颁发token
	token, err := jwt.GenToken(int64(user.ID), user.Name)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("颁发token失败")
	}

	resp.StatusCode = 0
	message := "登录成功"
	Id := int64(user.ID)
	resp.StatusMsg = &message
	resp.Token = &token
	resp.UserId = &Id
	return resp, nil
}

// GetUserInfo implements the BaseServiceImpl interface.
func (s *BaseServiceImpl) GetUserInfo(ctx context.Context, req *base.UserInfoRequest) (resp *base.UserInfoResponse, err error) {
	// TODO: Your code here...
	resp = base.NewUserInfoResponse()
	var message = ""
	resp.StatusMsg = &message

	var user *model.User
	if redis.UserIsExists(ctx, req.UserId) != 0 {
		user, err = redis.GetUserInfo(ctx, req.UserId)
		resp.StatusCode = 1
		if err != nil {
			return resp, err
		}
	} else {
		user, err = db.GetUserById(ctx, req.UserId)
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("数据库查询失败")
		}

		//	数据库查询成功后，需要将信息写入缓存
		redis.SetUserInfo(ctx, user)
	}

	isFollow := false
	if redis.FollowIsExists(ctx, req.UserId) != 0 {
		isFollow = redis.IsFollow(ctx, *req.ToUserId, req.UserId)
	} else {
		isFollow, err = db.IsFollowed(ctx, *req.ToUserId, req.UserId)
		if err != nil {
			return nil, err
		}
		var action int64
		if isFollow {
			action = 1
		} else {
			action = 0
		}
		redis.FollowAction(ctx, *req.ToUserId, req.UserId, action)
	}

	resp.StatusCode = 0
	message = "请求成功"
	resp.User = &base.UserInfo{
		Id:            int64(user.ID),
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      isFollow,
	}
	return resp, nil
}

// GetVideoList implements the BaseServiceImpl interface.
func (s *BaseServiceImpl) GetVideoList(ctx context.Context, req *base.FeedRequest) (resp *base.FeedResponse, err error) {
	// TODO: Your code here...
	resp = base.NewFeedResponse()
	message := ""
	nextTime := time.Now().Unix()
	resp.StatusMsg = &message
	resp.NextTime = &nextTime

	latestTime, err := strconv.Atoi(*req.LatestTime)
	if err != nil {
		return resp, errno.NewErrNo("转换数字时间失败")
	}
	fmt.Println("====================================================")
	fmt.Println("时间：", latestTime)
	fmt.Println("====================================================")
	videos, err := db.MGetByTime(ctx, int64(latestTime))
	if err != nil {
		resp.StatusCode = 1
		return nil, err
	}

	var videoList []*base.VideoInfo
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

		if redis.VideoIsExists(ctx, int64(video.ID)) == 0 {
			redis.SetVideoMessage(ctx, &video)
		}

		isLike := false
		if redis.LikeIsExists(ctx, *req.UserId) != 0 {
			isLike = redis.IsLike(ctx, *req.UserId, int64(video.ID))
		} else {
			isLike, err = db.HasLiked(ctx, *req.UserId, int64(video.ID))
			if err != nil {
				return
			}
			var action int64
			if isLike {
				action = 1
			} else {
				action = -1
			}
			redis.FavoriteAction(ctx, *req.UserId, int64(video.ID), action)
		}

		isFollow := false
		if redis.FollowIsExists(ctx, *req.UserId) != 0 {
			isFollow = redis.IsFollow(ctx, *req.UserId, video.UserId)
		} else {
			isFollow, _ = db.IsFollowed(ctx, *req.UserId, video.UserId)
			var action int64
			if isFollow {
				action = 1
			} else {
				action = -1
			}
			redis.FollowAction(ctx, *req.UserId, int64(video.ID), action)
		}

		video := &base.VideoInfo{
			Id: int64(video.ID),
			Author: &base.UserInfo{
				Id:            int64(user.ID),
				Name:          user.Name,
				FollowerCount: user.FollowerCount,
				IsFollow:      isFollow,
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
	message = "success"
	return resp, nil
}

// PublishAction implements the BaseServiceImpl interface.
func (s *BaseServiceImpl) PublishAction(ctx context.Context, request *base.PublishRequest) (resp *base.PublishResponse, err error) {
	// TODO: Your code here...
	resp = base.NewPublishResponse()
	message := ""
	resp.StatusMsg = &message
	VideoPath := request.Title + ".mp4"
	CoverPath := request.Title + ".jpg"

	//这里也需要保存到本地，路径那也是相对于当前的程序的
	if err = util.Cover("../../data/"+VideoPath, "../../data/"+CoverPath); err != nil {
		resp.StatusCode = 1
		str, _ := os.Getwd()
		return resp, errno.NewErrNo("获取视频封面失败，当前程序路径" + str)
	}
	// 执行数据库事务
	err = db.CreateVideo(ctx, VideoPath, CoverPath, request.Title, request.UserId)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("数据库寄寄")
	}

	resp.StatusCode = 0
	message = "success"
	return resp, nil
}

// GetPublishList implements the BaseServiceImpl interface.
func (s *BaseServiceImpl) GetPublishList(ctx context.Context, req *base.PublishListRequest) (resp *base.PublishListResponse, err error) {
	// TODO: Your code here...
	resp = new(base.PublishListResponse)
	message := ""
	resp.StatusMsg = &message

	var user *model.User
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
	fmt.Println("---------------------------------------------------------------")
	fmt.Println(*user)
	fmt.Println("---------------------------------------------------------------")
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("没找到视频相关信息，寄")
	}
	var videoList []*base.VideoInfo
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
					kv = append(kv, strconv.FormatInt(userLike, 10))
					kv = append(kv, "1")
				}
				if !redis.SetFavoriteList(ctx, req.UserId, kv...) {
					return resp, err
				}
			}
		}
		isLike = redis.IsLike(ctx, req.UserId, int64(video.ID))

		video := &base.VideoInfo{
			Id: int64(video.ID),
			Author: &base.UserInfo{
				Id:            int64(user.ID),
				Name:          user.Name,
				FollowerCount: user.FollowerCount,
				IsFollow:      false,
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
	message = "success"
	return resp, err
}

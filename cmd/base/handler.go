package main

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"simple_tiktok/cmd/base/dal/db"
	"simple_tiktok/cmd/base/dal/redis"
	base "simple_tiktok/kitex_gen/base"
	"simple_tiktok/util/errno"
	util "simple_tiktok/util/ffmpeg"
	"simple_tiktok/util/jwt"
	"strconv"
)

// BaseServiceImpl implements the last service interface defined in the IDL.
type BaseServiceImpl struct{}

// UserRegister implements the BaseServiceImpl interface.
func (s *BaseServiceImpl) UserRegister(ctx context.Context, request *base.RegisterRequest) (resp *base.RegisterResponse, err error) {
	// TODO: Your code here...
	resp = new(base.RegisterResponse)

	users, err := db.QueryUserByName(ctx, request.Username)
	if err != nil {
		resp.StatusCode = 1
		resp.StatusMsg = "数据库查询用户存不存在失败"
		return resp, err
	}
	if len(users) > 0 {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("用户已存在")
	}
	EncodePassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("密码加密失败")
	}
	newUser := &db.User{
		Username: request.Username,
		Password: string(EncodePassword),
	}
	id, err := db.CreateUser(ctx, newUser)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("创建用户失败")
	}

	token, err := jwt.GenToken(int64(newUser.ID), newUser.Username)
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
	users, err := db.QueryUserByName(ctx, request.Username)
	fmt.Println("==================================================")
	fmt.Println("这里是服务端", request)
	fmt.Println("==================================================")
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("数据库查询用户存不存在失败！")
	}

	fmt.Println("==================================================")
	fmt.Println("这里是服务端", resp)
	fmt.Println("==================================================")
	//密码认证
	err = bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(request.Password))
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("用户名或密码错误")
	}

	//颁发token
	token, err := jwt.GenToken(int64(users[0].ID), users[0].Username)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("颁发token失败")
	}

	resp.StatusCode = 0
	message := "登录成功"
	Id := int64(users[0].ID)
	resp.StatusMsg = &message
	resp.Token = &token
	resp.UserId = &Id
	return resp, nil
}

// GetUserInfo implements the BaseServiceImpl interface.
func (s *BaseServiceImpl) GetUserInfo(ctx context.Context, request *base.UserInfoRequest) (resp *base.UserInfoResponse, err error) {
	// TODO: Your code here...
	resp = new(base.UserInfoResponse)
	var message = ""
	resp.StatusMsg = &message

	fmt.Println("===================================")
	fmt.Println(request)
	fmt.Println("===================================")
	//数据库查一查，有没有这个人
	ThisUser, err := db.GetUser(ctx, request.UserId)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("数据库查询失败")
	}
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("没有此用户")
	}
	resp.StatusCode = 0
	message = "请求成功"
	resp.User = &base.UserInfo{
		Id:            int64(ThisUser.ID),
		Name:          ThisUser.Username,
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
	}
	return resp, nil
}

// GetVideoList implements the BaseServiceImpl interface.
func (s *BaseServiceImpl) GetVideoList(ctx context.Context, request *base.FeedRequest) (resp *base.FeedResponse, err error) {
	// TODO: Your code here...
	resp = new(base.FeedResponse)
	message := ""
	nextTime := int64(114514)
	resp.StatusMsg = &message
	resp.NextTime = &nextTime
	videos, err := db.FindVideoAll(ctx)
	if err != nil {
		message = "数据库查询失败"
		resp.StatusCode = 1
		return resp, errno.NewErrNo("数据库查询失败")
	}
	// 假设你有一个名为videoSlice的存储了多个VideoInfo的切片

	var videoList []*base.VideoInfo
	for _, v := range videos {
		//fmt.Println(v.CoverPath)
		users, _ := db.QueryUserByName(ctx, v.UserName)
		likeCount, _ := redis.GetLikeCount(ctx, int64(v.ID))

		// 这里以后还要改改
		commetCount := 0
		isLike := false
		likeInfo, _ := redis.GetLikeInfo(ctx, strconv.Itoa(int(v.ID)), strconv.Itoa(int(users[0].ID)))
		if likeInfo != nil {
			isLike = true
		}

		video := &base.VideoInfo{
			Id: int64(v.ID),
			Author: &base.UserInfo{
				Id:            int64(users[0].ID),
				Name:          v.UserName,
				FollowerCount: 0,
				IsFollow:      false,
			},
			PlayUrl:       "http://192.168.137.1:8888/data/" + v.VideoPath,
			CoverUrl:      "http://192.168.137.1:8888/data/" + v.CoverPath,
			FavoriteCount: likeCount,
			CommentCount:  int64(commetCount),
			IsFavorite:    isLike,
			Title:         v.Title,
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
	resp = new(base.PublishResponse)
	message := ""
	resp.StatusMsg = &message
	video := db.Video{
		UserName:  request.UserName,
		Title:     request.Title,
		VideoPath: request.Title + ".mp4",
		CoverPath: request.Title + ".jpg",
	}

	//这里也需要保存到本地，路径那也是相对于当前的程序的
	if err := util.Cover("../../data/"+video.VideoPath, "../../data/"+video.CoverPath); err != nil {
		resp.StatusCode = 1
		str, _ := os.Getwd()
		return resp, errno.NewErrNo("获取视频封面失败，当前程序路径" + str)
	}
	// 执行数据库事务
	if err := db.CreateVedio(ctx, &video); err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("数据库寄寄")
	}

	err = redis.InitLikeCount(ctx, int64(video.ID))
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("redis寄寄")
	}

	resp.StatusCode = 0
	message = "success"
	return resp, nil
}

// GetPublishList implements the BaseServiceImpl interface.
func (s *BaseServiceImpl) GetPublishList(ctx context.Context, request *base.PublishListRequest) (resp *base.PublishListResponse, err error) {
	// TODO: Your code here...
	resp = new(base.PublishListResponse)
	message := ""
	resp.StatusMsg = &message

	users, err := db.QueryUserById(ctx, request.UserId)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("查找这个用户失败")
	}

	videos, err := db.FindVideoByUserID(ctx, users[0].Username)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("没找到视频相关信息，寄")
	}
	var videoList []*base.VideoInfo
	for _, v := range videos {
		//fmt.Println(v.CoverPath)
		like, err := redis.GetLikeCount(ctx, int64(v.ID))
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("redis，寄")
		}
		likeInfo, err := redis.GetLikeInfo(ctx, strconv.Itoa(int(v.ID)), strconv.Itoa(int(users[0].ID)))
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("Redis查询是否点赞信息出错")
		}

		var isLike bool
		if likeInfo == nil {
			isLike = false
		} else {
			isLike = true
		}
		video := &base.VideoInfo{
			Id: int64(v.ID),
			Author: &base.UserInfo{
				Id:            int64(users[0].ID),
				Name:          users[0].Username,
				FollowerCount: 0,
				IsFollow:      false,
			},
			PlayUrl:       "http://192.168.137.1:8888/data/" + v.VideoPath,
			CoverUrl:      "http://192.168.137.1:8888/data/" + v.CoverPath,
			FavoriteCount: like,
			CommentCount:  0,
			IsFavorite:    isLike,
			Title:         v.Title,
		}
		videoList = append(videoList, video)
	}

	resp.VideoList = videoList
	message = "success"
	return resp, err
}

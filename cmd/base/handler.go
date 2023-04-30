package main

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"os"
	"simple_tiktok/dal/db"
	"simple_tiktok/dal/redis"
	base "simple_tiktok/kitex_gen/base"
	"simple_tiktok/pkg/errno"
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

	user, err := db.GetUserByName(ctx, request.Username)
	if err != nil {
		resp.StatusCode = 1
		resp.StatusMsg = "数据库查询用户存不存在失败"
		return resp, err
	}
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
func (s *BaseServiceImpl) GetUserInfo(ctx context.Context, request *base.UserInfoRequest) (resp *base.UserInfoResponse, err error) {
	// TODO: Your code here...
	resp = new(base.UserInfoResponse)
	var message = ""
	resp.StatusMsg = &message

	User, err := db.GetUserById(ctx, request.UserId)
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
		Id:            int64(User.ID),
		Name:          User.Name,
		FollowCount:   User.FollowCount,
		FollowerCount: User.FollowerCount,
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

	lastestTime, err := strconv.Atoi(*request.LatestTime)
	if err != nil {
		return resp, errno.NewErrNo("转换数字时间失败	")
	}

	// 数据库查询降序，Order(created_at desc)即可
	videos, err := db.MGetByTime(ctx, int64(lastestTime))
	if err != nil {
		resp.StatusCode = 1

	}
	// 假设你有一个名为videoSlice的存储了多个VideoInfo的切片

	var videoList []*base.VideoInfo
	for _, v := range videos {
		//fmt.Println(v.CoverPath)
		user, _ := db.GetUserById(ctx, v.UserId)
		likeCount, _ := redis.GetLikeCount(ctx, int64(v.ID))

		// 这里以后还要改改
		commetCount, err := db.QueryCommentsCount(ctx, int64(v.ID))
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("数据库查询是评论数量失败")
		}
		isLike := false
		isFollow := false

		if request.UserId != nil {
			isLike, err = redis.HasLiked(ctx, *request.UserId, strconv.Itoa(int(v.ID)))
			if err != nil {
				resp.StatusCode = 1
				return resp, errno.NewErrNo("redis查询是否点赞失败")
			}
			isFollow, err = redis.HasFollowed(ctx, *request.UserId, user.ID)
		}

		video := &base.VideoInfo{
			Id: int64(v.ID),
			Author: &base.UserInfo{
				Id:            int64(user.ID),
				Name:          user.Name,
				FollowerCount: 0,
				IsFollow:      isFollow,
			},
			PlayUrl:       "http://192.168.137.1:8888/data/" + v.PlayUrl,
			CoverUrl:      "http://192.168.137.1:8888/data/" + v.CoverUrl,
			FavoriteCount: likeCount,
			CommentCount:  commetCount,
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
func (s *BaseServiceImpl) GetPublishList(ctx context.Context, request *base.PublishListRequest) (resp *base.PublishListResponse, err error) {
	// TODO: Your code here...
	resp = new(base.PublishListResponse)
	message := ""
	resp.StatusMsg = &message

	user, err := db.GetUserById(ctx, request.UserId)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("查找这个用户失败")
	}

	videos, err := db.GetVideoByUserId(ctx, int64(user.ID))
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
		isLike, err := redis.HasLiked(ctx, int64(user.ID), strconv.Itoa(int(v.ID)))
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("Redis查询是否点赞信息出错")
		}

		video := &base.VideoInfo{
			Id: int64(v.ID),
			Author: &base.UserInfo{
				Id:            int64(user.ID),
				Name:          user.Name,
				FollowerCount: 0,
				IsFollow:      false,
			},
			PlayUrl:       "http://192.168.137.1:8888/data/" + v.PlayUrl,
			CoverUrl:      "http://192.168.137.1:8888/data/" + v.CoverUrl,
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

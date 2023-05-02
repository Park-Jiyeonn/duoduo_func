package main

import (
	"context"
	"fmt"
	"simple_tiktok/cmd/interact/mq"
	"simple_tiktok/dal/db"
	"simple_tiktok/dal/db/model"
	"simple_tiktok/dal/redis"
	"simple_tiktok/kitex_gen/base"
	interact "simple_tiktok/kitex_gen/interact"
	"simple_tiktok/pkg/errno"
	"strconv"
	"time"
)

// InteractServiceImpl implements the last service interface defined in the IDL.
type InteractServiceImpl struct{}

// LikeAction implements the InteractServiceImpl interface.
func (s *InteractServiceImpl) LikeAction(ctx context.Context, req *interact.LikeRequest) (resp *interact.LikeResponse, err error) {
	// TODO: Your code here...
	resp = interact.NewLikeResponse()

	if redis.LikeIsExists(ctx, *req.UserId) == 0 {
		likeList, err := db.GetUserLikeRecords(ctx, *req.UserId)
		if err != nil {
			resp.StatusCode = 1
			return resp, err
		}
		if len(likeList) > 0 {
			kv := make([]string, 0)
			for _, videoId := range likeList {
				kv = append(kv, strconv.FormatInt(*videoId, 10))
				kv = append(kv, "1")
			}
			if !redis.SetFavoriteList(ctx, *req.UserId, kv...) {
				resp.StatusCode = 1
				return resp, errno.NewErrNo("Redis设置用户点赞视频缓存出错")
			}
		}
	}
	if redis.UserIsExists(ctx, *req.UserId) == 0 {
		user, err := db.GetUserById(ctx, *req.UserId)
		if err != nil {
			resp.StatusCode = 1
			return resp, err
		}
		if !redis.SetUserInfo(ctx, user) {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("Redis缓存用户信息出错")
		}
	}
	if redis.LikeIsExists(ctx, req.VideoId) != 0 {
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
	if redis.UserIsExists(ctx, *req.UserId) == 0 {
		user, err := db.GetUserById(ctx, authorID)
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
	like := redis.IsLike(ctx, *req.UserId, req.VideoId)
	if like && req.GetActionType() == "2" {
		action = -1
	} else if !like && req.GetActionType() == "1" {
		action = 1
	} else {
		return
	}

	// 更新点赞列表
	if !redis.FavoriteAction(ctx, *req.UserId, req.VideoId, action) {
		return resp, errno.NewErrNo("更新点赞列表失败")
	}
	// 更新用户点赞数
	if !redis.IncrUserField(ctx, *req.UserId, "favorite_count", action) {
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

	resp.StatusCode = 0
	resp.StatusMsg = "success"
	return resp, nil
}

// GetLikeList implements the InteractServiceImpl interface.
func (s *InteractServiceImpl) GetLikeList(ctx context.Context, req *interact.LikeListRequest) (resp *interact.LikeListResponse, err error) {
	// TODO: Your code here...
	resp = interact.NewLikeListResponse()
	message := ""
	resp.StatusMsg = &message

	var likeList []*int64
	if redis.LikeIsExists(ctx, req.UserId) == 0 {
		likeList, err = db.GetUserLikeRecords(ctx, req.UserId)
		if err != nil {
			resp.StatusCode = 1
			return resp, err
		}
		if len(likeList) > 0 {
			kv := make([]string, 0)
			for _, videoId := range likeList {
				kv = append(kv, strconv.FormatInt(*videoId, 10))
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

	var videoList []*base.VideoInfo
	for _, vid := range likeList {
		var video *model.Video
		if redis.VideoIsExists(ctx, *vid) == 0 {
			video, err = db.GetVideoByVideoId(ctx, *vid)
			if err != nil {
				return nil, err
			}
			redis.SetVideoMessage(ctx, video)
		} else {
			video, err = redis.GetVideoMessage(ctx, *vid)
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

		retVideo := &base.VideoInfo{
			Id: *vid,
			Author: &base.UserInfo{
				Id:            video.UserId,
				Name:          user.Name,
				FollowerCount: 0,
				IsFollow:      false,
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
	message = "success"
	return resp, nil
}

// CommentAction implements the InteractServiceImpl interface.
func (s *InteractServiceImpl) CommentAction(ctx context.Context, request *interact.CommentRequest) (resp *interact.CommentResponse, err error) {
	// TODO: Your code here...
	resp = new(interact.CommentResponse)
	message := ""
	resp.StatusMsg = &message

	var newComment = model.Comment{
		UserId:      *request.UserId,
		VideoId:     request.VideoId,
		Content:     "",
		PublishDate: time.Now().Format("2006-01-02 15:04:05"),
	}
	if request.CommentText != nil {
		newComment.Content = *request.CommentText
	}
	if request.ActionType == "1" {
		err := mq.SendComment(&newComment)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("goroutine crashed:", err)
				}
			}()
			err := mq.ReceiveMessage(ctx)
			if err != nil {
				fmt.Println(err.Error())
			}
		}()
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("发送消息到消息队列失败！" + err.Error())
		}
	} else {
		err := db.DeleteCommentByID(ctx, *request.CommentId)
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("删除评论失败！")
		}
	}

	resp.StatusCode = 0
	message = "success"

	resComment := &interact.CommentInfo{
		Id: int64(newComment.ID),
		User: &base.UserInfo{
			Id:            *request.UserId,
			Name:          *request.UserName,
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      false,
		},
		Content:    newComment.Content,
		CreateDate: newComment.PublishDate,
	}
	message = "success"
	resp.Comment = resComment
	return resp, nil
}

// GetCommentList implements the InteractServiceImpl interface.
func (s *InteractServiceImpl) GetCommentList(ctx context.Context, request *interact.CommentListRequest) (resp *interact.CommentListResponse, err error) {
	// TODO: Your code here...
	resp = interact.NewCommentListResponse()
	message := ""
	resp.StatusMsg = &message
	comments, err := db.GetCommentList(ctx, request.VideoId)
	fmt.Println(comments)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("创建评论失败！")
	}
	var commentList []*interact.CommentInfo
	for _, v := range comments {
		user, err := db.GetUserById(ctx, v.UserId)
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("根据视频评论的ID查询用户失败！")
		}
		newComment := &interact.CommentInfo{
			Id: int64(v.ID),
			User: &base.UserInfo{
				Id:            int64(user.ID),
				Name:          user.Name,
				FollowCount:   0,
				FollowerCount: 0,
				IsFollow:      false,
			},
			Content:    v.Content,
			CreateDate: v.CreatedAt.String(),
		}
		commentList = append(commentList, newComment)
	}
	resp.CommentList = commentList
	message = "success"
	return resp, nil
}

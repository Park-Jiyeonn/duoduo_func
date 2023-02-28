package main

import (
	"context"
	"fmt"
	"simple_tiktok/cmd/interact/dal/db"
	"simple_tiktok/cmd/interact/dal/redis"
	"simple_tiktok/kitex_gen/base"
	interact "simple_tiktok/kitex_gen/interact"
	"simple_tiktok/util/errno"
	"strconv"
	"time"
)

// InteractServiceImpl implements the last service interface defined in the IDL.
type InteractServiceImpl struct{}

// LikeAction implements the InteractServiceImpl interface.
func (s *InteractServiceImpl) LikeAction(ctx context.Context, request *interact.LikeRequest) (resp *interact.LikeResponse, err error) {
	// TODO: Your code here...
	resp = new(interact.LikeResponse)
	likeInfo, err := redis.GetLikeInfo(ctx, request.VideoId, strconv.Itoa(int(*request.UserId)))
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("Redis查询是否点赞信息出错")
	}
	if likeInfo == nil {
		err := redis.SetLikeInfo(ctx, request.VideoId, strconv.Itoa(int(*request.UserId)))
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("Redis初始化点赞信息出错")
		}
		err = redis.IncrLikeCount(ctx, request.VideoId)
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("赞操作失败")
		}
	} else {
		err := redis.DecrLikeCount(ctx, request.VideoId)
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("取消赞操作失败")
		}
		err = redis.DelLikeInfo(ctx, request.VideoId, strconv.Itoa(int(*request.UserId)))
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("删除赞操作失败")
		}
	}
	resp.StatusCode = 0
	resp.StatusMsg = "success"
	return resp, nil
}

// GetLikeList implements the InteractServiceImpl interface.
func (s *InteractServiceImpl) GetLikeList(ctx context.Context, request *interact.LikeListRequest) (resp *interact.LikeListResponse, err error) {
	// TODO: Your code here...
	resp = new(interact.LikeListResponse)
	message := ""
	resp.StatusMsg = &message
	videos, err := redis.GetLikedVideos(ctx, request.UserId)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("redis查点赞过的列表，寄")
	}
	var videoList []*base.VideoInfo
	for _, v := range videos {
		//fmt.Println(v.CoverPath)
		like, err := redis.GetLikeCount(ctx, v)
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("redis，寄")
		}

		theVideo, err := db.FindVideoByVideoID(ctx, v)
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("没查到这个视频" + strconv.FormatInt(v, 10))
		}

		users, err := db.QueryUserByName(ctx, theVideo.UserName)
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("没查到这个作者" + theVideo.UserName)
		}

		video := &base.VideoInfo{
			Id: v,
			Author: &base.UserInfo{
				Id:            int64(users[0].ID),
				Name:          users[0].Username,
				FollowerCount: 0,
				IsFollow:      false,
			},
			PlayUrl:       "http://192.168.137.1:8888/data/" + theVideo.VideoPath,
			CoverUrl:      "http://192.168.137.1:8888/data/" + theVideo.CoverPath,
			FavoriteCount: like,
			CommentCount:  0,
			IsFavorite:    true,
			Title:         theVideo.Title,
		}
		videoList = append(videoList, video)
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

	var newComment = db.Comment{
		UserID:      *request.UserId,
		VideoID:     request.VideoId,
		Content:     "",
		PublishDate: time.Now().Format("2006-01-02 15:04:05"),
	}
	if request.CommentText != nil {
		newComment.Content = *request.CommentText
	}
	if request.ActionType == "1" {
		err := db.CreateComment(ctx, &newComment)
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("创建评论失败！")
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
	resp = new(interact.CommentListResponse)
	message := ""
	resp.StatusMsg = &message
	comments, err := db.QueryCommentByVideoID(ctx, request.VideoId)
	fmt.Println(comments)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("创建评论失败！")
	}
	var commentList []*interact.CommentInfo
	for _, v := range comments {
		users, err := db.QueryUserById(ctx, v.UserID)
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("根据视频评论的ID查询用户失败！")
		}
		newComment := &interact.CommentInfo{
			Id: int64(v.ID),
			User: &base.UserInfo{
				Id:            int64(users[0].ID),
				Name:          users[0].Username,
				FollowCount:   0,
				FollowerCount: 0,
				IsFollow:      false,
			},
			Content:    v.Content,
			CreateDate: v.PublishDate,
		}
		commentList = append(commentList, newComment)
	}
	resp.CommentList = commentList
	message = "success"
	return resp, nil
}

package handler

import (
	"context"
	api_model "duoduo_fun/api/model"
	"duoduo_fun/dal/db"
	"duoduo_fun/dal/db/model"
	"duoduo_fun/dal/redis"
	"duoduo_fun/pkg/errno"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func commentAction(ctx context.Context, request *api_model.CommentRequest) (resp *api_model.CommentResponse, err error) {
	var newComment = model.Comment{
		UserId:      request.UserId,
		VideoId:     request.VideoId,
		Content:     "",
		PublishDate: time.Now().Format("2006-01-02 15:04:05"),
	}

	if request.ActionType == "1" {
		if err := db.CreateComment(ctx, &newComment); err != nil {
			return nil, err
		}
		redis.IncrVideoField(ctx, newComment.VideoId, "comment_count", 1)
		if err != nil {
			resp.StatusCode = 1
			return nil, err
		}
	} else {
		err := db.DeleteCommentByID(ctx, request.CommentId)
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("删除评论失败！")
		}
		redis.IncrVideoField(ctx, request.CommentId, "comment_id", -1)
	}

	resp.StatusCode = 200
	resp.StatusMsg = "success"

	resComment := api_model.CommentInfo{
		Id: int(newComment.ID),
		User: api_model.UserInfo{
			Id:   request.UserId,
			Name: request.UserName,
		},
		Content:    newComment.Content,
		CreateDate: newComment.PublishDate,
	}
	resp.Comment = resComment
	return resp, nil
}
func CommentAction(c *gin.Context) {
	req := &api_model.CommentRequest{}
	err := c.BindJSON(req)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := commentAction(context.Background(), req)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func getCommentList(ctx context.Context, request *api_model.CommentListRequest) (resp *api_model.CommentListResponse, err error) {
	comments, err := db.GetCommentList(ctx, request.VideoId)
	fmt.Println(comments)
	if err != nil {
		resp.StatusCode = 1
		return resp, errno.NewErrNo("查询评论失败！")
	}
	var commentList []api_model.CommentInfo
	for _, v := range comments {
		user, err := db.GetUserById(ctx, v.UserId)
		if err != nil {
			resp.StatusCode = 1
			return resp, errno.NewErrNo("根据视频评论的ID查询用户失败！")
		}
		newComment := api_model.CommentInfo{
			Id: int(v.ID),
			User: api_model.UserInfo{
				Id:   int(user.ID),
				Name: user.Name,
			},
			Content:    v.Content,
			CreateDate: v.CreatedAt.String(),
		}
		commentList = append(commentList, newComment)
	}
	resp.CommentList = commentList
	resp.StatusMsg = "success"
	return resp, nil
}
func GetCommentList(c *gin.Context) {
	req := &api_model.CommentListRequest{}
	err := c.BindJSON(req)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := getCommentList(context.Background(), req)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

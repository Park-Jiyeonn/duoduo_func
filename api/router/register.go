package router

import (
	"duoduo_fun/api/handler"
	"duoduo_fun/api/mw"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	api := r.Group("/api")

	treasure := api.Group("/treasure")
	{
		treasure.POST("", handler.PrizeDraw)
	}

	user := api.Group("/user")
	{
		user.POST("/register", handler.UserRegister)
		user.POST("/login", handler.UserLogin)
		user.GET("/info", handler.GetUserInfo)
	}

	video := api.Group("/video")
	{
		video.GET("/feed", handler.GetVideoList)
		video.POST("/publish", handler.PublishAction)
		video.GET("/list", handler.GetPublishList)
	}

	like := api.Group("/like")
	{
		like.POST("", handler.LikeAction).Use(mw.IPLimitMiddleware())
		like.GET("/list", handler.GetLikeList)
	}

	comment := api.Group("/comment")
	{
		comment.POST("", handler.CommentAction).Use(mw.TokenLimitMiddleware())
		comment.GET("/list", handler.GetCommentList)
	}
}

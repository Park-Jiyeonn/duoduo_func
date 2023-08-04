package router

import (
	"duoduo_fun/api/handler"
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
	}
}

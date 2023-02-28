package mw

import (
	"context"
	"fmt"
	"simple_tiktok/util/jwt"

	"github.com/cloudwego/hertz/pkg/app"
)

func JwtMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		auth := c.Query("token")
		// URL中为检测到token
		if auth == "" {
			auth = c.PostForm("token")
		}

		mc, err := jwt.ParseToken(auth)
		if err != nil {
			fmt.Println("输出下token看看")
			fmt.Println("token = ", auth)
			c.Abort()
			return
		}

		// 将当前请求的username信息保存到请求的上下文c上
		// 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
		fmt.Println("鉴权成功")
		fmt.Println(mc)
		c.Set("userid", mc.UserID)
		c.Set("username", mc.Username)
		c.Next(ctx)
	}
}

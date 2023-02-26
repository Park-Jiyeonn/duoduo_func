package mw

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"simple_tiktok/util/jwt"
)

func JwtMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// pre-handle
		// ...
		fmt.Println("中间件被使用了吗？")
		token := c.Query("token")
		if token == "" {
			token = c.PostForm("token")
			// 未携带Token
			if token == "" {
				c.Abort()
			}
		}
		fmt.Println("token = ", token)
		claim, err := jwt.ParseToken(token)
		if err != nil || claim == nil {
			fmt.Println("这个token没搞对")
			c.Abort()
		} else {
			c.Set("user_name", claim.Username)
			fmt.Println(claim)
			//fmt.Println(token)
			c.Next(ctx)
		}
	}
}

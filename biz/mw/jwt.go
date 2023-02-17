package mw

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"simple_tiktok/util"
)

func JwtMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// pre-handle
		// ...
		fmt.Println("中间件被使用了吗？")
		token := c.Query("token")
		if token == "" {
			c.Abort()
		}
		claim, err := util.ParseToken(token)
		if err != nil {
			c.Abort()
		}
		c.Set("user_name", claim.Username)
		c.Next(ctx)
	}
}

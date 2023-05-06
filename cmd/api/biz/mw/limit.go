package mw

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"simple_tiktok/dal/redis"
	"simple_tiktok/pkg/consts"
	"time"
)

func IPLimitMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		clientIp := c.ClientIP()
		key := consts.GetFavoriteLmtKey(clientIp)
		fmt.Println("==================================================")
		fmt.Println(clientIp)
		fmt.Println("key = ", key)
		fmt.Println("==================================================")
		if redis.Exists(ctx, key) == 0 {
			// redis中不存在当前key
			fmt.Println("我在这里！！！！！！！！！！！！！！！！！！！！！！！！！！！！！")
			redis.Set(ctx, key, 1, time.Minute)
		} else {
			cnt, _ := redis.Incr(ctx, key)
			if cnt >= consts.Limits_per_sec {
				// 操作过于频繁
				c.String(http.StatusBadRequest, "操作过于频繁")
				c.Abort()
				return
			}
		}

		c.Next(ctx)
	}
}

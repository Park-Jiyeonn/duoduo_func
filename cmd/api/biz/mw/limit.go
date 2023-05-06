package mw

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"simple_tiktok/dal/redis"
	"simple_tiktok/pkg/consts"
	"simple_tiktok/pkg/errno"
	"time"
)

func IPLimitMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		clientIp := c.ClientIP()
		key := consts.GetFavoriteLmtKey(clientIp)
		if redis.Exists(ctx, key) == 0 {
			// redis中不存在当前key
			redis.Set(ctx, key, 1, time.Minute)
		} else {
			cnt, _ := redis.Incr(ctx, key)
			if cnt >= consts.Limits_per_sec {
				// 操作过于频繁
				c.JSON(http.StatusBadRequest, errno.NewErrNo("操作过于频繁！请稍后再试"))
				c.Abort()
				return
			}
		}

		c.Next(ctx)
	}
}

package mw

import (
	"context"
	"duoduo_fun/dal/redis"
	"duoduo_fun/pkg/consts"
	"duoduo_fun/pkg/errno"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

func IPLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIp := c.ClientIP()
		ctx := context.Background()
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
		c.Next()
	}
}

var tokenLimiter *rate.Limiter

func TokenLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if tokenLimiter == nil {
			tokenLimiter = rate.NewLimiter(rate.Every(time.Millisecond), 100)
		}
		if tokenLimiter.Allow() {
			c.String(http.StatusBadRequest, "当前请求太多，请重试")
			c.Abort()
			return
		}
		c.Next()
	}
}

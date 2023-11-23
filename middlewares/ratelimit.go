package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) { //令牌桶:与漏桶相比,一瞬间产生大量请求
	bucket := ratelimit.NewBucket(fillInterval, cap) // ,如果cap // 容量里的令牌数够多，会迅速全部处理
	return func(c *gin.Context) {
		// 如果取不到令牌就中断本次请求返回 rate limit...
		if bucket.TakeAvailable(1) == 0 {
			c.String(http.StatusOK, "rate limit...")
			c.Abort()
			return
		}
		//取到令牌就放行
		c.Next()
	}
}

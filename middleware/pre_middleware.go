package middleware

import (
	"NilCTF/error_code"
	"NilCTF/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// PreMiddleware 包含限速器的中间件处理
type PreMiddleware struct {
	shardLimiter *utils.ShardLimiter
}

// NewPreMiddleware 初始化 Handler
func NewPreMiddleware() *PreMiddleware {
	return &PreMiddleware{}
}

// RateLimitMiddleware 创建基于 IP 的限速中间件
func (h *PreMiddleware) RateLimitMiddleware(r rate.Limit, b int, cacheSize int) gin.HandlerFunc {
	h.shardLimiter = utils.NewShardLimiter(r, b, cacheSize)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := h.shardLimiter.GetLimiter(ip)
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status":  "fail",
				"message": error_code.ErrTooManyRequests.Error(),
			})
			return
		}
		c.Next()
	}
}

// CSPMiddleware 设置内容安全策略头
func (h *PreMiddleware) CSPMiddleware(cspValue string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Security-Policy", cspValue)
		c.Next()
	}
}

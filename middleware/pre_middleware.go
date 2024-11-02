package middleware

import (
	"NilCTF/error_code"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/time/rate"
)

// ipLimiter 用于维护 IP 地址到限速器的映射
type ipLimiter struct {
	limiters sync.Map // 使用 sync.Map 提高并发性能
	rate     rate.Limit
	burst    int
}

// newIPLimiter 初始化 ipLimiter 并启动定期清理过期限速器的 goroutine
func newIPLimiter(r rate.Limit, b int) *ipLimiter {
	i := &ipLimiter{
		rate:  r,
		burst: b,
	}

	// 启动一个goroutine定期清理过期的限速器
	go func() {
		for {
			time.Sleep(1 * time.Minute) // 每分钟清理一次
			i.cleanupLimiters()
		}
	}()

	return i
}

// getLimiter 返回指定 IP 的限速器，如果不存在则创建新的
func (i *ipLimiter) getLimiter(ip string) *rate.Limiter {
	if limiter, exists := i.limiters.Load(ip); exists {
		return limiter.(*rate.Limiter)
	}

	// 创建一个新的限速器
	newLimiter := rate.NewLimiter(i.rate, i.burst)
	i.limiters.Store(ip, newLimiter)
	return newLimiter
}

// cleanupLimiters 清理未活跃的限速器
func (i *ipLimiter) cleanupLimiters() {
	i.limiters.Range(func(ip, limiter interface{}) bool {
		lim := limiter.(*rate.Limiter)
		// 如果限速器在5分钟内没有被使用过，则清理它
		if lim.AllowN(time.Now(), 0) {
			i.limiters.Delete(ip)
		}
		return true
	})
}

// Handler 包含所有中间件相关的处理器
type PreMiddleware struct {
	ipLimiter *ipLimiter
}

// NewPreMiddleware 初始化 Handler
func NewPreMiddleware() *PreMiddleware {
	return &PreMiddleware{}
}

// RateLimitMiddleware 是基于 IP 的速率限制中间件
func (h *PreMiddleware) RateLimitMiddleware(r rate.Limit, b int) gin.HandlerFunc {
	h.ipLimiter = newIPLimiter(r, b)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := h.ipLimiter.getLimiter(ip)

		// 检查限速器是否允许请求
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
func (h *PreMiddleware) CSPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline';")
		c.Next()
	}
}

// BluemondayMiddleware 过滤表单输入
func (h *PreMiddleware) BluemondayMiddleware(maxParamCount, maxKeyLength, maxFieldLength int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析表单数据
		if err := c.Request.ParseForm(); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": error_code.ErrFailedToParseForm.Error(),
			})
			return
		}

		// 检查参数数量
		if len(c.Request.PostForm) > maxParamCount {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": error_code.ErrTooManyParameters.Error(),
			})
			return
		}

		// 清理每个输入字段并检查字符长度
		for key, values := range c.Request.PostForm {
			if len(key) > maxKeyLength {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"status":  "fail",
					"message": error_code.ErrKeyTooLong.Error(),
				})
				return
			}
			for i, value := range values {
				if len(value) > maxFieldLength {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
						"status":  "fail",
						"message": error_code.ErrInputTooLong.Error(),
					})
					return
				}

				p := bluemonday.UGCPolicy()
				c.Request.PostForm[key][i] = p.Sanitize(value) // 清理输入并替换原值
			}
		}

		c.Next()
	}
}

// LimitRequestBody 通过参数 maxBytes 动态限制请求体的大小
func (h *PreMiddleware) LimitRequestBody(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查请求体大小
		if c.Request.ContentLength > maxBytes {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"status": "fail", "message": error_code.ErrRequestBodyTooLarge.Error()})
			return
		}
		c.Next()
	}
}

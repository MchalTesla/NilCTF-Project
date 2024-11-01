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
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	rate     rate.Limit
	burst    int
}

func newIPLimiter(r rate.Limit, b int) *ipLimiter {
	return &ipLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}
}

func (i *ipLimiter) getLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	// 如果限速器已经存在，直接返回
	if limiter, exists := i.limiters[ip]; exists {
		return limiter
	}

	// 否则，创建一个新的限速器
	limiter := rate.NewLimiter(i.rate, i.burst)
	i.limiters[ip] = limiter

	// 启动一个 goroutine 定期清理过期的限速器（例如，5 分钟后）
	go func() {
		time.Sleep(5 * time.Minute)
		i.mu.Lock()
		delete(i.limiters, ip)
		i.mu.Unlock()
	}()

	return limiter
}

// Handler 包含所有中间件相关的处理器
type Handler struct {
	ipLimiter *ipLimiter
}

// NewHandler 初始化 Handler 及其限速器
func NewHandler(r rate.Limit, b int) *Handler {
	return &Handler{
		ipLimiter: newIPLimiter(r, b),
	}
}

// RateLimitMiddleware 是基于 IP 的速率限制中间件
func (h *Handler) RateLimitMiddleware() gin.HandlerFunc {
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
func (h *Handler) CSPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline';")
		c.Next()
	}
}

// BluemondayMiddleware 过滤表单输入
func (h *Handler) BluemondayMiddleware(maxParamCount, maxKeyLength, maxFieldLength int) gin.HandlerFunc {
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

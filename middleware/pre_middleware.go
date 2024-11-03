package middleware

import (
	"NilCTF/error_code"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/time/rate"
)

// shardLimiter 使用 LRU 缓存存储分片限速器
type shardLimiter struct {
	cache      *lru.Cache
	mu         sync.Mutex
	rate       rate.Limit
	burst      int
	shardCount int
}

// 创建一个带 LRU 缓存的分片限速器管理器
func newShardLimiter(rate rate.Limit, burst int, cacheSize int) *shardLimiter {
	cache, err := lru.New(cacheSize)
	if err != nil {
		return nil
	}
	return &shardLimiter{
		cache:      cache,
		rate:       rate,
		burst:      burst,
		shardCount: cacheSize * 2, // 根据缓存大小（最大用户数）* 2 设置分片数量
	}
}

// getLimiter 获取或创建限速器
func (sl *shardLimiter) getLimiter(ip string) *rate.Limiter {
	shardKey := hashIP(ip) % uint32(sl.shardCount)

	// 先检查缓存
	if limiter, ok := sl.cache.Get(shardKey); ok {
		return limiter.(*rate.Limiter)
	}

	// 使用双重检查获取限速器
	sl.mu.Lock()
	defer sl.mu.Unlock()

	if limiter, ok := sl.cache.Get(shardKey); ok {
		return limiter.(*rate.Limiter)
	}

	// 创建新限速器并存入缓存
	limiter := rate.NewLimiter(sl.rate, sl.burst)
	sl.cache.Add(shardKey, limiter)
	return limiter
}

// hashIP 哈希函数确保分片均匀
func hashIP(ip string) uint32 {
	var hash uint32 = 2166136261
	for i := 0; i < len(ip); i++ {
		hash ^= uint32(ip[i])
		hash *= 16777619
	}
	return hash
}

// PreMiddleware 包含限速器的中间件处理
type PreMiddleware struct {
	shardLimiter *shardLimiter
}

// NewPreMiddleware 初始化 Handler
func NewPreMiddleware() *PreMiddleware {
	return &PreMiddleware{}
}

// RateLimitMiddleware 创建基于 IP 的限速中间件
func (h *PreMiddleware) RateLimitMiddleware(r rate.Limit, b int, cacheSize int) gin.HandlerFunc {
	h.shardLimiter = newShardLimiter(r, b, cacheSize)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := h.shardLimiter.getLimiter(ip)
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
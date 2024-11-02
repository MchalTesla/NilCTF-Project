package middleware

import (
	"NilCTF/error_code"
	"encoding/json"
	"net/http"
	"sync"
	"unicode/utf8"
	"html"

	"bytes"
	"github.com/gin-gonic/gin"
	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/time/rate"
	"io"
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

// EscapeStringMiddleware 处理参数过滤和内容检查
func (h *PreMiddleware) FilterRequestParameters(maxParamCount, maxKeyLength, maxFieldLength int, maxFileSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodGet:
			if err := h.processGETRequest(c, maxParamCount, maxKeyLength, maxFieldLength); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"status":  "fail",
					"message": err.Error(),
				})
			}
		case http.MethodPost:
			if err := h.processPOSTRequest(c, maxParamCount, maxKeyLength, maxFieldLength, maxFileSize); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"status":  "fail",
					"message": err.Error(),
				})
			}
		default:
			c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{
				"status":  "fail",
				"message": error_code.ErrUnsupportedContentType.Error(),
			})
		}
		c.Next()
	}
}

// processGETRequest 处理 GET 请求
func (h *PreMiddleware) processGETRequest(c *gin.Context, maxParamCount, maxKeyLength, maxFieldLength int) error {
	queryParams := c.Request.URL.Query()

	if len(queryParams) > maxParamCount {
		return error_code.ErrTooManyParameters
	}

	escapedParams := make(map[string][]string)
	for key, values := range queryParams {
		if utf8.RuneCountInString(key) > maxKeyLength {
			return error_code.ErrKeyTooLong
		}
		escapedKey := html.EscapeString(key)
		for _, value := range values {
			if utf8.RuneCountInString(value) > maxFieldLength {
				return error_code.ErrInputTooLong
			}
			escapedParams[escapedKey] = append(escapedParams[escapedKey], html.EscapeString(value))
		}
	}
	c.Request.URL.RawQuery = queryToString(escapedParams)
	return nil
}

// processPOSTRequest 处理 POST 请求
func (h *PreMiddleware) processPOSTRequest(c *gin.Context, maxParamCount, maxKeyLength, maxFieldLength int, maxFileSize int64) error {
	switch c.ContentType() {
	case "application/json":
		var jsonData map[string]interface{}
		if err := c.ShouldBindJSON(&jsonData); err != nil {
			return error_code.ErrInternalServer
		}
		if len(jsonData) > maxParamCount {
			return error_code.ErrTooManyParameters
		}
		sanitizeJSONData(jsonData, maxKeyLength, maxFieldLength)
		jsonBytes, err := json.Marshal(jsonData)
		if err != nil {
			return error_code.ErrInternalServer
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(jsonBytes))
		c.Request.ContentLength = int64(len(jsonBytes))
	case "multipart/form-data":
		if err := c.Request.ParseMultipartForm(maxFileSize); err != nil {
			return error_code.ErrFileTooLarge
		}
		if len(c.Request.MultipartForm.File) > maxParamCount {
			return error_code.ErrTooManyFiles
		}
		for _, files := range c.Request.MultipartForm.File {
			for _, file := range files {
				if file.Size > maxFileSize {
					return error_code.ErrFileTooLarge
				}
			}
		}
	default:
		return error_code.ErrUnsupportedContentType
	}
	return nil
}

// sanitizeJSONData 递归处理 JSON 数据的键和值
func sanitizeJSONData(jsonData map[string]interface{}, maxKeyLength, maxFieldLength int) {
	for key, value := range jsonData {
		if utf8.RuneCountInString(key) > maxKeyLength {
			delete(jsonData, key)
			continue
		}
		escapedKey := html.EscapeString(key)

		switch v := value.(type) {
		case string:
			jsonData[escapedKey] = html.EscapeString(v)
		case []interface{}:
			for i, item := range v {
				if strItem, ok := item.(string); ok {
					v[i] = html.EscapeString(strItem)
				} else if itemMap, ok := item.(map[string]interface{}); ok {
					sanitizeJSONData(itemMap, maxKeyLength, maxFieldLength)
				}
			}
		case map[string]interface{}:
			sanitizeJSONData(v, maxKeyLength, maxFieldLength)
		}
		if key != escapedKey {
			delete(jsonData, key)
		}
	}
}

// queryToString 将查询参数转换为字符串
func queryToString(params map[string][]string) string {
	query := ""
	for key, values := range params {
		for _, value := range values {
			query += key + "=" + value + "&"
		}
	}
	if len(query) == 0 {
        return query // 如果 query 为空，则直接返回
    }
	return query[:len(query)-1]
}
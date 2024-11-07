package utils

import (
	"NilCTF/error_code"
	"crypto/rand"
	"time"
	"sync"

	"golang.org/x/time/rate"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hashicorp/golang-lru"
)

// PreMiddle
// ShardLimiter 使用 LRU 缓存存储分片限速器
type ShardLimiter struct {
	cache		*lru.Cache
	mu			sync.Mutex
	rate		rate.Limit
	burst		int
	shardCount	int
}

// 创建一个带 LRU 缓存的分片限速器管理器
func NewShardLimiter(rate rate.Limit, burst int, cacheSize int) *ShardLimiter {
	cache, err := lru.New(cacheSize)
	if err != nil {
		return nil
	}
	return &ShardLimiter{
		cache:      cache,
		rate:       rate,
		burst:      burst,
		shardCount: cacheSize * 2, // 根据缓存大小（最大用户数）* 2 设置分片数量
	}
}

// getLimiter 获取或创建限速器
func (sl *ShardLimiter) GetLimiter(ip string) *rate.Limiter {
	shardKey := HashIP(ip) % uint32(sl.shardCount)

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
func HashIP(ip string) uint32 {
	var hash uint32 = 2166136261
	for i := 0; i < len(ip); i++ {
		hash ^= uint32(ip[i])
		hash *= 16777619
	}
	return hash
}


// PostMiddle
// 生成随机字节
func GenerateRandomSecret(length int) ([]byte, error) {
	secret := make([]byte, length)
	_, err := rand.Read(secret) // 生成随机字节
	if err != nil {
		return nil, err
	}
	return secret, nil
}

// Claims 自定义声明结构体，继承 jwt.RegisteredClaims
type Claims struct {
	ID uint `json:"ID"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func GenerateToken(ID uint, jwtTime int, jwtSecret []byte) (string, error) {
	claims := Claims{
		ID: ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtTime) * time.Hour)), // 设置过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // 设置签发时间
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", error_code.ErrInternalServer // 生成错误转换为内部服务器错误
	}
	return tokenString, nil
}

// parseToken 解析 JWT Token 并返回声明
func ParseToken(tokenString string, jwtSecret []byte) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		if err == jwt.ErrTokenExpired {
			return nil, nil, error_code.ErrTokenExpired
		}
		return nil, nil, error_code.ErrInternalServer // 解析错误转换为内部服务器错误
	}

	return token, claims, nil
}
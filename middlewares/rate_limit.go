package middlewares

import (
	"net/http"
	"sync"
	"time"

	"gitee.com/NextEraAbyss/gin-template/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter 基于IP的限流器.
type IPRateLimiter struct {
	ips    map[string]*rate.Limiter
	mu     *sync.RWMutex
	r      rate.Limit
	b      int
	expire time.Duration
	// 记录每个IP最后访问时间
	lastSeen map[string]time.Time
}

// NewIPRateLimiter 创建一个新的IP限流器.
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips:      make(map[string]*rate.Limiter),
		mu:       &sync.RWMutex{},
		r:        r,
		b:        b,
		expire:   time.Hour, // 默认1小时后过期
		lastSeen: make(map[string]time.Time),
	}

	// 启动一个协程定期清理过期的限流器
	go i.cleanupExpired()
	return i
}

// AddIP 为指定IP创建一个新的速率限制器并添加到映射中
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = limiter
	i.lastSeen[ip] = time.Now()
	return limiter
}

// GetLimiter 返回指定IP的速率限制器，如果不存在则创建一个
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	limiter, exists := i.ips[ip]
	i.mu.RUnlock()

	if !exists {
		return i.AddIP(ip)
	}

	// 更新最后访问时间
	i.mu.Lock()
	i.lastSeen[ip] = time.Now()
	i.mu.Unlock()

	return limiter
}

// cleanupExpired 定期清理过期的IP限制器
func (i *IPRateLimiter) cleanupExpired() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		i.mu.Lock()
		for ip, lastTime := range i.lastSeen {
			// 如果在一段时间内没有请求，则删除该IP的限制器
			if now.Sub(lastTime) > i.expire {
				delete(i.ips, ip)
				delete(i.lastSeen, ip)
			}
		}
		i.mu.Unlock()
	}
}

// RateLimit 创建一个限流中间件
// 参数:
// - max: 在time.Duration内允许的最大请求数
// - duration: 限流的时间窗口
func RateLimit(max int, duration time.Duration) gin.HandlerFunc {
	// 计算速率限制 (请求/秒)
	limit := rate.Limit(float64(max) / duration.Seconds())
	limiter := NewIPRateLimiter(limit, max)

	return func(c *gin.Context) {
		// 获取客户端IP
		ip := c.ClientIP()
		// 获取该IP的限流器
		ipLimiter := limiter.GetLimiter(ip)
		// 尝试获取令牌
		if !ipLimiter.Allow() {
			// 如果无法获取令牌，返回429状态码
			utils.ResponseError(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}
		c.Next()
	}
}

package middleware

import (
	"context"
	"sync"
	"time"

	"github.com/YOJIA-yukino/simple-douyin-backend/internal/config"
	"github.com/YOJIA-yukino/simple-douyin-backend/internal/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"golang.org/x/time/rate"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(jwtManager *utils.JWTManager) app.HandlerFunc {
	return func(ctx *app.RequestContext) {
		token := ctx.Query("token")
		if token == "" {
			// 尝试从Header获取
			token = ctx.GetHeader("Authorization")
			if token != "" && len(token) > 7 {
				token = token[7:] // 移除 "Bearer " 前缀
			}
		}

		if token == "" {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"status_code": 401,
				"status_msg":  "token is required",
			})
			ctx.Abort()
			return
		}

		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"status_code": 401,
				"status_msg":  "invalid token",
			})
			ctx.Abort()
			return
		}

		// 将用户信息存储到上下文中
		ctx.Set("user_id", claims.UserID)
		ctx.Set("username", claims.Username)
		ctx.Next()
	}
}

// RateLimitMiddleware 限流中间件
type RateLimitMiddleware struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimitMiddleware 创建限流中间件
func NewRateLimitMiddleware(rps int, burst int) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(rps),
		burst:    burst,
	}
}

// getLimiter 获取限流器
func (m *RateLimitMiddleware) getLimiter(key string) *rate.Limiter {
	m.mu.Lock()
	defer m.mu.Unlock()

	limiter, exists := m.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(m.rate, m.burst)
		m.limiters[key] = limiter
	}

	return limiter
}

// RateLimit 限流处理
func (m *RateLimitMiddleware) RateLimit() app.HandlerFunc {
	return func(ctx *app.RequestContext) {
		// 使用IP作为限流key
		key := ctx.ClientIP()
		limiter := m.getLimiter(key)

		if !limiter.Allow() {
			ctx.JSON(consts.StatusTooManyRequests, map[string]interface{}{
				"status_code": 429,
				"status_msg":  "rate limit exceeded",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// LoggingMiddleware 日志中间件
func LoggingMiddleware() app.HandlerFunc {
	return func(ctx *app.RequestContext) {
		start := time.Now()
		path := string(ctx.Request.URI().Path())
		method := string(ctx.Request.Method())

		// 记录请求开始
		hlog.Infof("Request started: %s %s", method, path)

		ctx.Next()

		// 记录请求结束
		duration := time.Since(start)
		status := ctx.Response.StatusCode()
		hlog.Infof("Request completed: %s %s - %d - %v", method, path, status, duration)
	}
}

// CORSMiddleware CORS中间件
func CORSMiddleware(cfg *config.Config) app.HandlerFunc {
	return func(ctx *app.RequestContext) {
		ctx.Header("Access-Control-Allow-Origin", cfg.Security.AllowedOrigins)
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if string(ctx.Request.Method()) == "OPTIONS" {
			ctx.AbortWithStatus(consts.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware() app.HandlerFunc {
	return func(ctx *app.RequestContext) {
		defer func() {
			if err := recover(); err != nil {
				hlog.Errorf("Panic recovered: %v", err)
				ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
					"status_code": 500,
					"status_msg":  "internal server error",
				})
				ctx.Abort()
			}
		}()

		ctx.Next()
	}
}

// TimeoutMiddleware 超时中间件
func TimeoutMiddleware(timeout time.Duration) app.HandlerFunc {
	return func(ctx *app.RequestContext) {
		// 创建带超时的上下文
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		// 将超时上下文设置到请求上下文中
		ctx.Set("timeout_ctx", timeoutCtx)

		// 监听超时
		go func() {
			select {
			case <-timeoutCtx.Done():
				if timeoutCtx.Err() == context.DeadlineExceeded {
					hlog.Warnf("Request timeout: %s", string(ctx.Request.URI().Path()))
					ctx.JSON(consts.StatusRequestTimeout, map[string]interface{}{
						"status_code": 408,
						"status_msg":  "request timeout",
					})
					ctx.Abort()
				}
			}
		}()

		ctx.Next()
	}
}

// MetricsMiddleware 指标中间件
type MetricsMiddleware struct {
	requestCount int64
	errorCount   int64
	responseTime time.Duration
	mu           sync.RWMutex
}

// NewMetricsMiddleware 创建指标中间件
func NewMetricsMiddleware() *MetricsMiddleware {
	return &MetricsMiddleware{}
}

// Metrics 指标处理
func (m *MetricsMiddleware) Metrics() app.HandlerFunc {
	return func(ctx *app.RequestContext) {
		start := time.Now()

		ctx.Next()

		duration := time.Since(start)
		status := ctx.Response.StatusCode()

		m.mu.Lock()
		m.requestCount++
		if status >= 400 {
			m.errorCount++
		}
		m.responseTime = duration
		m.mu.Unlock()

		// 记录指标
		hlog.Infof("Metrics: requests=%d, errors=%d, response_time=%v, status=%d",
			m.requestCount, m.errorCount, duration, status)
	}
}

// GetMetrics 获取指标
func (m *MetricsMiddleware) GetMetrics() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]interface{}{
		"request_count": m.requestCount,
		"error_count":   m.errorCount,
		"response_time": m.responseTime.String(),
		"error_rate":    float64(m.errorCount) / float64(m.requestCount),
	}
}

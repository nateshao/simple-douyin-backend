package monitoring

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/YOJIA-yukino/simple-douyin-backend/internal/cache"
	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
)

// HealthChecker 健康检查器
type HealthChecker struct {
	db    *gorm.DB
	cache cache.Cache
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker(db *gorm.DB, cache cache.Cache) *HealthChecker {
	return &HealthChecker{
		db:    db,
		cache: cache,
	}
}

// HealthStatus 健康状态
type HealthStatus struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Services  map[string]ServiceInfo `json:"services"`
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Latency string `json:"latency,omitempty"`
}

// CheckHealth 检查健康状态
func (h *HealthChecker) CheckHealth(ctx context.Context) *HealthStatus {
	status := &HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services:  make(map[string]ServiceInfo),
	}

	// 检查数据库
	dbInfo := h.checkDatabase(ctx)
	status.Services["database"] = dbInfo
	if dbInfo.Status != "healthy" {
		status.Status = "unhealthy"
	}

	// 检查Redis
	redisInfo := h.checkRedis(ctx)
	status.Services["redis"] = redisInfo
	if redisInfo.Status != "healthy" {
		status.Status = "unhealthy"
	}

	return status
}

// checkDatabase 检查数据库健康状态
func (h *HealthChecker) checkDatabase(ctx context.Context) ServiceInfo {
	start := time.Now()

	// 执行简单查询测试连接
	var result int
	err := h.db.WithContext(ctx).Raw("SELECT 1").Scan(&result).Error

	latency := time.Since(start)

	if err != nil {
		return ServiceInfo{
			Status:  "unhealthy",
			Message: fmt.Sprintf("Database connection failed: %v", err),
		}
	}

	return ServiceInfo{
		Status:  "healthy",
		Message: "Database connection is working",
		Latency: latency.String(),
	}
}

// checkRedis 检查Redis健康状态
func (h *HealthChecker) checkRedis(ctx context.Context) ServiceInfo {
	start := time.Now()

	// 执行PING命令测试连接
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := h.cache.Exists(ctx, "health_check")
	latency := time.Since(start)

	if err != nil {
		return ServiceInfo{
			Status:  "unhealthy",
			Message: fmt.Sprintf("Redis connection failed: %v", err),
		}
	}

	return ServiceInfo{
		Status:  "healthy",
		Message: "Redis connection is working",
		Latency: latency.String(),
	}
}

// Metrics 指标收集器
type Metrics struct {
	RequestCount      int64
	ErrorCount        int64
	ResponseTime      time.Duration
	ActiveConnections int64
	mu                sync.RWMutex
}

// NewMetrics 创建指标收集器
func NewMetrics() *Metrics {
	return &Metrics{}
}

// IncrementRequest 增加请求计数
func (m *Metrics) IncrementRequest() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RequestCount++
}

// IncrementError 增加错误计数
func (m *Metrics) IncrementError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ErrorCount++
}

// UpdateResponseTime 更新响应时间
func (m *Metrics) UpdateResponseTime(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ResponseTime = duration
}

// UpdateActiveConnections 更新活跃连接数
func (m *Metrics) UpdateActiveConnections(count int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ActiveConnections = count
}

// GetMetrics 获取指标
func (m *Metrics) GetMetrics() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	errorRate := float64(0)
	if m.RequestCount > 0 {
		errorRate = float64(m.ErrorCount) / float64(m.RequestCount)
	}

	return map[string]interface{}{
		"request_count":      m.RequestCount,
		"error_count":        m.ErrorCount,
		"response_time":      m.ResponseTime.String(),
		"active_connections": m.ActiveConnections,
		"error_rate":         errorRate,
		"success_rate":       1 - errorRate,
	}
}

// 全局变量
var (
	healthChecker *HealthChecker
	metrics       *Metrics
)

// InitMonitoring 初始化监控
func InitMonitoring(db *gorm.DB, cache cache.Cache) {
	healthChecker = NewHealthChecker(db, cache)
	metrics = NewMetrics()
}

// HealthCheckHandler 健康检查处理器
func HealthCheckHandler(ctx context.Context, c *app.RequestContext) {
	if healthChecker == nil {
		c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status":  "unhealthy",
			"message": "Health checker not initialized",
		})
		return
	}

	status := healthChecker.CheckHealth(ctx)

	if status.Status == "healthy" {
		c.JSON(http.StatusOK, status)
	} else {
		c.JSON(http.StatusServiceUnavailable, status)
	}
}

// MetricsHandler 指标处理器
func MetricsHandler(ctx context.Context, c *app.RequestContext) {
	if metrics == nil {
		c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status":  "error",
			"message": "Metrics not initialized",
		})
		return
	}

	c.JSON(http.StatusOK, metrics.GetMetrics())
}

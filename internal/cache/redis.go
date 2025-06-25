package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/YOJIA-yukino/simple-douyin-backend/internal/config"
	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis(cfg *config.Config) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	return nil
}

// GetRedisClient 获取Redis客户端
func GetRedisClient() *redis.Client {
	return redisClient
}

// Close 关闭Redis连接
func Close() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}

// SetWithExpiration 设置带过期时间的键值对
func SetWithExpiration(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return redisClient.Set(ctx, key, value, expiration).Err()
}

// Get 获取值
func Get(ctx context.Context, key string) (string, error) {
	return redisClient.Get(ctx, key).Result()
}

// Del 删除键
func Del(ctx context.Context, keys ...string) error {
	return redisClient.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	result, err := redisClient.Exists(ctx, key).Result()
	return result > 0, err
}

// HMSet 设置哈希表字段
func HMSet(ctx context.Context, key string, values map[string]interface{}) error {
	return redisClient.HMSet(ctx, key, values).Err()
}

// HMGet 获取哈希表字段
func HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	return redisClient.HMGet(ctx, key, fields...).Result()
}

// Expire 设置过期时间
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return redisClient.Expire(ctx, key, expiration).Err()
}

// Incr 递增
func Incr(ctx context.Context, key string) error {
	return redisClient.Incr(ctx, key).Err()
}

// Decr 递减
func Decr(ctx context.Context, key string) error {
	return redisClient.Decr(ctx, key).Err()
}

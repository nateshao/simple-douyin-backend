package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client
var redisCache *RedisCache

// RedisCache 适配器，实现 Cache 接口
type RedisCache struct {
	client *redis.Client
}

// InitRedis 初始化Redis连接
func InitRedis(addr string, password string, db int, poolSize int, minIdleConns int) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	})
	redisCache = &RedisCache{client: redisClient}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}
	return nil
}

// GetRedisClient 获取Redis缓存适配器
func GetRedisClient() Cache {
	return redisCache
}

// SetWithExpiration 设置带过期时间的键值对
func (r *RedisCache) SetWithExpiration(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get 获取值
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Del 删除键
func (r *RedisCache) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	return result > 0, err
}

// HMSet 设置哈希表字段
func (r *RedisCache) HMSet(ctx context.Context, key string, values map[string]interface{}) error {
	return r.client.HMSet(ctx, key, values).Err()
}

// HMGet 获取哈希表字段
func (r *RedisCache) HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	return r.client.HMGet(ctx, key, fields...).Result()
}

// Expire 设置过期时间
func (r *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

// Incr 递增
func (r *RedisCache) Incr(ctx context.Context, key string) error {
	return r.client.Incr(ctx, key).Err()
}

// Decr 递减
func (r *RedisCache) Decr(ctx context.Context, key string) error {
	return r.client.Decr(ctx, key).Err()
}

// Close 关闭Redis连接
func (r *RedisCache) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

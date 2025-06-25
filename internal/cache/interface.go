package cache

import (
	"context"
	"time"
)

// Cache 缓存接口
type Cache interface {
	// 基本操作
	SetWithExpiration(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)

	// 哈希表操作
	HMSet(ctx context.Context, key string, values map[string]interface{}) error
	HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error)

	// 过期时间
	Expire(ctx context.Context, key string, expiration time.Duration) error

	// 计数器操作
	Incr(ctx context.Context, key string) error
	Decr(ctx context.Context, key string) error

	// 列表操作
	LPush(ctx context.Context, key string, values ...interface{}) error
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	LRem(ctx context.Context, key string, count int64, value interface{}) error
	LTrim(ctx context.Context, key string, start, stop int64) error

	// 关闭连接
	Close() error
}

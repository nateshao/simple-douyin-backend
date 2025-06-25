package service

import (
	"math/rand"
	"sync"
	"time"

	"github.com/YOJIA-yukino/simple-douyin-backend/internal/cache"
)

var (
	redisClient cache.Cache
	redisOnce   sync.Once
)

func initRedis() {
	redisOnce.Do(func() {
		redisClient = cache.GetRedisClient()
	})
}

const (
	emptyCache           = "{}"
	emptyCacheExpireTime = time.Hour
)

func getEmptyCacheExpireTime() time.Duration {
	return time.Duration(int64(emptyCacheExpireTime) + rand.Int63n(int64(30*time.Minute)))
}

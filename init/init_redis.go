package init

import (
	"fmt"

	"github.com/YOJIA-yukino/simple-douyin-backend/internal/cache"
	"github.com/go-redis/redis"
)

var rdb *redis.Client

func InitRDB() {
	// 使用新的缓存接口初始化Redis
	err := cache.InitRedis(
		fmt.Sprintf("%s:%s", rdbHost, rdbPort),
		"",  // password
		0,   // db
		100, // poolSize
		10,  // minIdleConns
	)

	if err != nil {
		stdOutLogger.Error().Caller().Str("Redis启动失败", err.Error()).Msg("Redis initialization failed, continuing without Redis")
		return
	}

	stdOutLogger.Info().Msg("Redis initialized successfully")
}

func GetRDB() *redis.Client {
	return rdb
}

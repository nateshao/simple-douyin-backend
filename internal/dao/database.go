package dao

import (
	"sync"

	initialization "github.com/YOJIA-yukino/simple-douyin-backend/init"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
	cache  cache.Cache
)

// DaoInitialization 初始化Dao层的服务，包括获取DB以及Kafka
func DaoInitialization() {
	dbOnce.Do(func() {
		db = initialization.GetDB()
		initKafkaClient()

		// 初始化缓存
		cache = cache.GetRedisClient()
	})
}

// GetCache 获取缓存实例
func GetCache() cache.Cache {
	return cache
}

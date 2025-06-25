package dao

import (
	"sync"

	initialization "github.com/YOJIA-yukino/simple-douyin-backend/init"
	"github.com/YOJIA-yukino/simple-douyin-backend/internal/cache"
	"gorm.io/gorm"
)

var (
	db        *gorm.DB
	dbOnce    sync.Once
	cacheInst cache.Cache
)

// DaoInitialization 初始化Dao层的服务，包括获取DB以及Kafka
func DaoInitialization() {
	dbOnce.Do(func() {
		db = initialization.GetDB()
		initKafkaClient()

		// 初始化缓存
		cacheInst = cache.GetRedisClient()
	})
}

// GetCache 获取缓存实例
func GetCache() cache.Cache {
	return cacheInst
}

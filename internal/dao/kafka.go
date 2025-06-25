package dao

import (
	"sync"

	"github.com/Shopify/sarama"
	initialization "github.com/YOJIA-yukino/simple-douyin-backend/init"
	"github.com/YOJIA-yukino/simple-douyin-backend/internal/utils/logger"
)

var (
	kafkaClient sarama.Consumer
	kafkaOnce   sync.Once
)

func initKafkaClient() {
	kafkaOnce.Do(func() {
		kafkaClient = initialization.GetKafkaClient()
		// 检查Kafka客户端是否成功初始化
		if kafkaClient == nil {
			logger.GlobalLogger.Printf("Warning: Kafka client is nil, skipping message queue initialization")
			return
		}

		go func() {
			for {
				err := GetFavoriteDaoInstance().getFromMessageQueue()
				if err == nil {
					break
				}
				logger.GlobalLogger.Printf("Error in message queue processing: %v", err)
			}
		}()
	})
}

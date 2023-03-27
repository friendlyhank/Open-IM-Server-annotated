package msg_transfer

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/kafka"
)

var (
	persistentCH PersistentConsumerHandler
	producer     *kafka.Producer
)

func Init() {
	persistentCH.Init() // 持久化消息到数据库
	producer = kafka.NewKafkaProducer(config.Config.Kafka.Ms2pschat.Addr, config.Config.Kafka.Ms2pschat.Topic)
}

package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/kafka"
)

const ConsumerMsgs = 3 // 消费消息指令
const ChannelNum = 100 // 设置channel数量

var (
	persistentCH PersistentConsumerHandler         // 持久化数据mysql
	historyCH    OnlineHistoryRedisConsumerHandler //
	producer     *kafka.Producer
)

func Init() {
	persistentCH.Init() // 持久化消息到数据库
	historyCH.Init()
	producer = kafka.NewKafkaProducer(config.Config.Kafka.Ms2pschat.Addr, config.Config.Kafka.Ms2pschat.Topic)
}

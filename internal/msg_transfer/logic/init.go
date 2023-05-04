package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/kafka"
	"fmt"
)

const ConsumerMsgs = 3        // 消费消息指令
const AggregationMessages = 4 // 聚合消息
const ChannelNum = 100        // 设置channel数量

var (
	persistentCH PersistentConsumerHandler         // 持久化数据mysql
	historyCH    OnlineHistoryRedisConsumerHandler //
	producer     *kafka.Producer
)

func Init() {
	persistentCH.Init() // 持久化消息到数据库
	historyCH.Init()
	// 推送消息的生产者
	producer = kafka.NewKafkaProducer(config.Config.Kafka.Ms2pschat.Addr, config.Config.Kafka.Ms2pschat.Topic)
}

func Run() {
	if config.Config.ChatPersistenceMysql {
		go persistentCH.persistentConsumerGroup.RegisterHandleAndConsumer(&persistentCH)
	} else {
		fmt.Println("not start mysql consumer")
	}
	go historyCH.historyConsumerGroup.RegisterHandleAndConsumer(&historyCH)
}

package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/kafka"
)

/*
 * 初始化推送
 */

var (
	producer *kafka.Producer
)

func Init(rpcPort int) {

}

func init() {
	producer = kafka.NewKafkaProducer(config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.Ws2mschat.Topic)
}

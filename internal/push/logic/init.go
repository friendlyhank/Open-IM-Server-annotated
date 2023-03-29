package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/kafka"
)

/*
 * 初始化推送
 */

var (
	rpcServer RPCServer
	producer  *kafka.Producer
)

func Init(rpcPort int) {
	rpcServer.Init(rpcPort)
}

func init() {
	producer = kafka.NewKafkaProducer(config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.Ws2mschat.Topic)
}

func Run() {
	go rpcServer.run()
}

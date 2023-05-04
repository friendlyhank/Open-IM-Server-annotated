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
	pushCh    PushConsumerHandler
	producer  *kafka.Producer
)

func Init(rpcPort int) {
	rpcServer.Init(rpcPort) // 推送的rpc服务
	pushCh.Init()           // 消费kafka推送消息
}

func init() {
	producer = kafka.NewKafkaProducer(config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.Ws2mschat.Topic)
}

func Run() {
	go rpcServer.run()
	go pushCh.pushConsumerGroup.RegisterHandleAndConsumer(&pushCh) // 注册推送的消费者组
}

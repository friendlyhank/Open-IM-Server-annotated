package kafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
)

/*
 * 消费者组
 */

type MConsumerGroup struct {
	sarama.ConsumerGroup
	groupID string
	topics  []string
}

type MConsumerGroupConfig struct {
	KafkaVersion sarama.KafkaVersion // kafka版本
}

func NewMConsumerGroup(consumerConfig *MConsumerGroupConfig, topics, addrs []string, groupID string) *MConsumerGroup {
	config := sarama.NewConfig()
	config.Version = consumerConfig.KafkaVersion
	consumerGroup, err := sarama.NewConsumerGroup(addrs, groupID, config)
	if err != nil {
		fmt.Println("args:", addrs, groupID, config)
		panic(err.Error())
	}
	return &MConsumerGroup{
		consumerGroup,
		groupID,
		topics,
	}
}

// RegisterHandleAndConsumer - 注册消费者组
func (mc *MConsumerGroup) RegisterHandleAndConsumer(handler sarama.ConsumerGroupHandler) {
	ctx := context.Background()
	for {
		err := mc.ConsumerGroup.Consume(ctx, mc.topics, handler)
		if err != nil {
			panic(err.Error())
		}
	}
}

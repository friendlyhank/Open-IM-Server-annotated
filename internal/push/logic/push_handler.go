package logic

import (
	"Open_IM/pkg/common/config"
	kfk "Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	pbChat "Open_IM/pkg/proto/msg"
	pbPush "Open_IM/pkg/proto/push"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

/*
 * 消费kafka推送的数据
 */

type fcb func(msg []byte)

type PushConsumerHandler struct {
	msgHandle         map[string]fcb
	pushConsumerGroup *kfk.MConsumerGroup
}

func (ms *PushConsumerHandler) Init() {
	ms.msgHandle = make(map[string]fcb)
	ms.msgHandle[config.Config.Kafka.Ms2pschat.Topic] = ms.handleMs2PsChat
	ms.pushConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V2_0_0_0},
		[]string{config.Config.Kafka.Ms2pschat.Topic}, config.Config.Kafka.Ms2pschat.Addr,
		config.Config.Kafka.ConsumerGroupID.MsgToPush)
}

// handleMs2PsChat - 处理消费推送消息
func (ms *PushConsumerHandler) handleMs2PsChat(msg []byte) {
	log.NewDebug("", "msg come from kafka  And push!!!", "msg", string(msg))
	msgFromMQ := pbChat.PushMsgDataToMQ{}
	if err := proto.Unmarshal(msg, &msgFromMQ); err != nil {
		log.Error("", "push Unmarshal msg err", "msg", string(msg), "err", err.Error())
		return
	}
	pbData := &pbPush.PushMsgReq{
		OperationID:  msgFromMQ.OperationID,
		MsgData:      msgFromMQ.MsgData,
		PushToUserID: msgFromMQ.PushToUserID,
	}
	switch msgFromMQ.MsgData.SessionType {
	default:
		MsgToUser(pbData)
	}
}

func (PushConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (PushConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (ms *PushConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.NewDebug("", "kafka get info to mysql", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value))
		ms.msgHandle[msg.Topic](msg.Value)
		sess.MarkMessage(msg, "")
	}
	return nil
}

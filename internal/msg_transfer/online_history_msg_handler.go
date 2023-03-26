package msg_transfer

import (
	"Open_IM/pkg/common/config"
	kfk "Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	pbMsg "Open_IM/pkg/proto/msg"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

type fcb func(cMsg *sarama.ConsumerMessage, msgKey string, sess sarama.ConsumerGroupSession)

type OnlineHistoryRedisConsumerHandler struct {
	msgHandle            map[string]fcb
	historyConsumerGroup *kfk.MConsumerGroup
}

func (och *OnlineHistoryRedisConsumerHandler) Init() {
	och.msgHandle[config.Config.Kafka.Ws2mschat.Topic] = och.handleChatWs2MongoLowReliability
	och.historyConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V2_0_0_0},
		[]string{config.Config.Kafka.Ws2mschat.Topic},
		config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.ConsumerGroupID.MsgToRedis)
}

func (och *OnlineHistoryRedisConsumerHandler) handleChatWs2MongoLowReliability(cMsg *sarama.ConsumerMessage, msgKey string, sess sarama.ConsumerGroupSession) {
	msg := cMsg.Value
	msgFromMQ := pbMsg.MsgDataToMQ{}
	err := proto.Unmarshal(msg, &msgFromMQ)
	if err != nil {
		log.Error("msg_transfer Unmarshal msg err", "", "msg", string(msg), "err", err.Error())
		return
	}
	operationID := msgFromMQ.OperationID
	log.NewInfo(operationID, "msg come mongo!!!", "", "msg", string(msg))
	//Control whether to store offline messages (mongo)
}

func (OnlineHistoryRedisConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (OnlineHistoryRedisConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (och *OnlineHistoryRedisConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error { // a instance in the consumer group
	cMsg := make([]*sarama.ConsumerMessage, 0, 1000)
	for msg := range claim.Messages() {
		if len(msg.Value) != 0 {
			cMsg = append(cMsg, msg)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

package msg_transfer

/*
 *  消费kafka消息并持久化到mysql
 */

import (
	"Open_IM/pkg/common/config"
	kfk "Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	pbMsg "Open_IM/pkg/proto/msg"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

type PersistentConsumerHandler struct {
	msgHandle               map[string]fcb
	persistentConsumerGroup *kfk.MConsumerGroup
}

func (pc *PersistentConsumerHandler) Init() {
	pc.msgHandle = make(map[string]fcb)
	pc.msgHandle[config.Config.Kafka.Ws2mschat.Topic] = pc.handleChatWs2Mysql
	pc.persistentConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V2_0_0_0},
		[]string{config.Config.Kafka.Ws2mschat.Topic},
		config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.ConsumerGroupID.MsgToMySql)
}

// handleChatWs2Mysql - 处理消费消息并持久化
func (pc *PersistentConsumerHandler) handleChatWs2Mysql(cMsg *sarama.ConsumerMessage, msgKey string, _ sarama.ConsumerGroupSession) {
	msg := cMsg.Value
	log.NewInfo("msg come here mysql!!!", "", "msg", string(msg), msgKey)
	msgFromMQ := pbMsg.MsgDataToMQ{}
	err := proto.Unmarshal(msg, &msgFromMQ)
	if err != nil {
		log.NewError(msgFromMQ.OperationID, "msg_transfer Unmarshal msg err", "msg", string(msg), "err", err.Error())
		return
	}
	log.Debug(msgFromMQ.OperationID, "proto.Unmarshal MsgDataToMQ", msgFromMQ.String())
	//Control whether to store history messages (mysql)

	log.NewInfo(msgFromMQ.OperationID, "msg_transfer msg persisting", string(msg))
	if err = im_mysql_msg_model.InsertMessageToChatLog(msgFromMQ); err != nil {
		log.NewError(msgFromMQ.OperationID, "Message insert failed", "err", err.Error(), "msg", msgFromMQ.String())
		return
	}
}

func (PersistentConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (PersistentConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (pc *PersistentConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.NewDebug("", "kafka get info to mysql", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value), "key", string(msg.Key))
		if len(msg.Value) != 0 {
			pc.msgHandle[msg.Topic](msg, string(msg.Key), sess)
		} else {
			log.Error("", "msg get from kafka but is nil", msg.Key)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

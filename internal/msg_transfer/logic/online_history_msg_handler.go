package logic

import (
	"Open_IM/pkg/common/config"
	kfk "Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbMsg "Open_IM/pkg/proto/msg"
	pbPush "Open_IM/pkg/proto/push"
	"context"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"strings"
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

	go sendMessageToPush(&msgFromMQ, msgKey)
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

// sendMessageToPush - 用推送的方式发送消息
func sendMessageToPush(message *pbMsg.MsgDataToMQ, pushToUserID string) {
	log.Info(message.OperationID, "msg_transfer send message to push", "message", message.String())
	rpcPushMsg := pbPush.PushMsgReq{OperationID: message.OperationID, MsgData: message.MsgData, PushToUserID: pushToUserID}
	mqPushMsg := pbMsg.PushMsgDataToMQ{OperationID: message.OperationID, MsgData: message.MsgData, PushToUserID: pushToUserID}
	grpcConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImPushName, message.OperationID)
	if grpcConn != nil {
		log.Error(rpcPushMsg.OperationID, "rpc dial failed", "push data", rpcPushMsg.String())
		pid, offset, err := producer.SendMessage(&mqPushMsg, mqPushMsg.PushToUserID, rpcPushMsg.OperationID)
		if err != nil {
			log.Error(mqPushMsg.OperationID, "kafka send failed", "send data", message.String(), "pid", pid, "offset", offset, "err", err.Error())
		}
		return
	}
	msgClient := pbPush.NewPushMsgServiceClient(grpcConn)
	// 推送消息
	_, err := msgClient.PushMsg(context.Background(), &rpcPushMsg)
	if err != nil {
		log.Error(rpcPushMsg.OperationID, "rpc send failed", rpcPushMsg.OperationID, "push data", rpcPushMsg.String(), "err", err.Error())
		pid, offset, err := producer.SendMessage(&mqPushMsg, mqPushMsg.PushToUserID, rpcPushMsg.OperationID)
		if err != nil {
			log.Error(message.OperationID, "kafka send failed", mqPushMsg.OperationID, "send data", mqPushMsg.String(), "pid", pid, "offset", offset, "err", err.Error())
		}
	} else {
		log.Info(message.OperationID, "rpc send success", rpcPushMsg.OperationID, "push data", rpcPushMsg.String())

	}
}

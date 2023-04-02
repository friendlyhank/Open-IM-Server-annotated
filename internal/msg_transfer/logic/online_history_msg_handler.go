package logic

import (
	"Open_IM/pkg/common/config"
	kfk "Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbMsg "Open_IM/pkg/proto/msg"
	pbPush "Open_IM/pkg/proto/push"
	"Open_IM/pkg/utils"
	"context"
	"github.com/Shopify/sarama"
	"strings"
	"sync"
	"time"
)

type fcb func(cMsg *sarama.ConsumerMessage, msgKey string, sess sarama.ConsumerGroupSession)

// 处理指令和对应的值
type Cmd2Value struct {
	Cmd   int
	Value interface{}
}

// 目标channel值
type TriggerChannelValue struct {
	triggerID string
	cmsgList  []*sarama.ConsumerMessage
}

type OnlineHistoryRedisConsumerHandler struct {
	msgHandle            map[string]fcb
	historyConsumerGroup *kfk.MConsumerGroup
	chArrays             [ChannelNum]chan Cmd2Value
	msgDistributionCh    chan Cmd2Value // 分发的channel
}

func (och *OnlineHistoryRedisConsumerHandler) Init() {
	och.historyConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V2_0_0_0},
		[]string{config.Config.Kafka.Ws2mschat.Topic},
		config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.ConsumerGroupID.MsgToRedis)
}

func (OnlineHistoryRedisConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (OnlineHistoryRedisConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (och *OnlineHistoryRedisConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error { // a instance in the consumer group
	rwLock := new(sync.RWMutex)
	cMsg := make([]*sarama.ConsumerMessage, 0, 1000)
	log.NewDebug("", "online new session msg come", claim.HighWaterMarkOffset(), claim.Topic(), claim.Partition())
	t := time.NewTicker(time.Duration(100) * time.Millisecond)
	var triggerID string
	go func() {
		for {
			select {
			case <-t.C:
				if len(cMsg) > 0 {
					rwLock.Lock()
					ccMsg := make([]*sarama.ConsumerMessage, 0, 1000)
					for _, v := range cMsg {
						ccMsg = append(ccMsg, v)
					}
					cMsg = make([]*sarama.ConsumerMessage, 0, 1000)
					rwLock.Unlock()
					split := 1000
					triggerID = utils.OperationIDGenerator()
					log.Debug(triggerID, "timer trigger msg consumer start", len(ccMsg))
					// 每1000条消息放入分发的channel中
					for i := 0; i < len(ccMsg)/split; i++ {
						//log.Debug()
						och.msgDistributionCh <- Cmd2Value{Cmd: ConsumerMsgs, Value: TriggerChannelValue{
							triggerID: triggerID, cmsgList: ccMsg[i*split : (i+1)*split]}}
					}
					if (len(ccMsg) % split) > 0 {
						och.msgDistributionCh <- Cmd2Value{Cmd: ConsumerMsgs, Value: TriggerChannelValue{
							triggerID: triggerID, cmsgList: ccMsg[split*(len(ccMsg)/split):]}}
					}
					//sess.MarkMessage(ccMsg[len(cMsg)-1], "")

					log.Debug(triggerID, "timer trigger msg consumer end", len(cMsg))
				}
			}
		}
	}()
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

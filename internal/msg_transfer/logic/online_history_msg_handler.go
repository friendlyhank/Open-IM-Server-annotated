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
	"github.com/golang/protobuf/proto"
	"hash/crc32"
	"strings"
	"sync"
	"time"
)

// MsgChannelValue - 聚合消息结构体
type MsgChannelValue struct {
	aggregationID string //maybe userID or super groupID
	triggerID     string
	msgList       []*pbMsg.MsgDataToMQ
	lastSeq       uint64
}

// TriggerChannelValue - 消费消息指令结构体
type TriggerChannelValue struct {
	triggerID string
	cmsgList  []*sarama.ConsumerMessage
}

type fcb func(cMsg *sarama.ConsumerMessage, msgKey string, sess sarama.ConsumerGroupSession)

// 处理指令和对应的值
type Cmd2Value struct {
	Cmd   int
	Value interface{}
}

type OnlineHistoryRedisConsumerHandler struct {
	msgHandle            map[string]fcb
	historyConsumerGroup *kfk.MConsumerGroup
	chArrays             [ChannelNum]chan Cmd2Value // 消息聚合指令chan
	msgDistributionCh    chan Cmd2Value             // 消息分发指令chan
}

func (och *OnlineHistoryRedisConsumerHandler) Init() {
	och.msgDistributionCh = make(chan Cmd2Value) //no buffer channel
	go och.MessagesDistributionHandle()          // 分发消息到聚合的chan
	for i := 0; i < ChannelNum; i++ {
		och.chArrays[i] = make(chan Cmd2Value, 50)
		go och.Run(i)
	}
	och.historyConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V2_0_0_0},
		[]string{config.Config.Kafka.Ws2mschat.Topic},
		config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.ConsumerGroupID.MsgToRedis)
}

// Run - 聚合后的消息处理
func (och *OnlineHistoryRedisConsumerHandler) Run(channelID int) {
	for {
		select {
		case cmd := <-och.chArrays[channelID]:
			switch cmd.Cmd {
			case AggregationMessages:
				msgChannelValue := cmd.Value.(MsgChannelValue)
				msgList := msgChannelValue.msgList
				triggerID := msgChannelValue.triggerID
				// todo hank 搞清楚storage和noStorage的区别
				storageMsgList := make([]*pbMsg.MsgDataToMQ, 0, 80)
				notStoragePushMsgList := make([]*pbMsg.MsgDataToMQ, 0, 80)
				log.Debug(triggerID, "msg arrived channel", "channel id", channelID, msgList, msgChannelValue.aggregationID, len(msgList))
				for _, v := range msgList {
					log.Debug(triggerID, "msg come to storage center", v.String())
					storageMsgList = append(storageMsgList, v)
				}
				log.Debug(triggerID, "msg storage length", len(storageMsgList), "push length", len(notStoragePushMsgList))
				if len(storageMsgList) > 0 {
					for _, v := range storageMsgList {
						sendMessageToPushMQ(v, msgChannelValue.aggregationID)
					}
				} else {

				}
			}
		}
	}
}

// MessagesDistributionHandle - 消息分发处理
func (och *OnlineHistoryRedisConsumerHandler) MessagesDistributionHandle() {
	for {
		aggregationMsgs := make(map[string][]*pbMsg.MsgDataToMQ, ChannelNum)
		select {
		case cmd := <-och.msgDistributionCh:
			switch cmd.Cmd {
			case ConsumerMsgs: // 消费消息的指令
				triggerChannelValue := cmd.Value.(TriggerChannelValue)
				triggerID := triggerChannelValue.triggerID       // 链路追踪的唯一id
				consumerMessages := triggerChannelValue.cmsgList // 消费的消息列表
				//Aggregation map[userid]message list
				log.Debug(triggerID, "batch messages come to distribution center", len(consumerMessages))
				for i := 0; i < len(consumerMessages); i++ {
					msgFromMQ := pbMsg.MsgDataToMQ{}
					err := proto.Unmarshal(consumerMessages[i].Value, &msgFromMQ)
					if err != nil {
						log.Error(triggerID, "msg_transfer Unmarshal msg err", "msg", string(consumerMessages[i].Value), "err", err.Error())
						return
					}
					log.Debug(triggerID, "single msg come to distribution center", msgFromMQ.String(), string(consumerMessages[i].Key))
					// todo hank why
					// 根据sendid或receiveid汇总消息
					if oldM, ok := aggregationMsgs[string(consumerMessages[i].Key)]; ok {
						oldM = append(oldM, &msgFromMQ)
						aggregationMsgs[string(consumerMessages[i].Key)] = oldM
					} else {
						m := make([]*pbMsg.MsgDataToMQ, 0, 100)
						m = append(m, &msgFromMQ)
						aggregationMsgs[string(consumerMessages[i].Key)] = m
					}
				}
				log.Debug(triggerID, "generate map list users len", len(aggregationMsgs))
				for aggregationID, v := range aggregationMsgs {
					if len(v) >= 0 {
						hashCode := getHashCode(aggregationID)
						channelID := hashCode % ChannelNum
						log.Debug(triggerID, "generate channelID", hashCode, channelID, aggregationID)
						//go func(cID uint32, userID string, messages []*pbMsg.MsgDataToMQ) {
						och.chArrays[channelID] <- Cmd2Value{Cmd: AggregationMessages, Value: MsgChannelValue{aggregationID: aggregationID, msgList: v, triggerID: triggerID}}
					}
				}
			}
		}
	}
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

	// 处理消息到分发消费消息的指令中
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
					//todo hank 为什么这样设计 每1000条消息放入分发的channel中
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
		rwLock.Lock()
		if len(msg.Value) != 0 {
			cMsg = append(cMsg, msg)
		}
		rwLock.Unlock()
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

// sendMessageToPushMQ - 发送推送的消息到mq
func sendMessageToPushMQ(message *pbMsg.MsgDataToMQ, pushToUserID string) {
	log.Info(message.OperationID, "msg_transfer send message to push", "message", message.String())
	rpcPushMsg := pbPush.PushMsgReq{OperationID: message.OperationID, MsgData: message.MsgData, PushToUserID: pushToUserID}
	mqPushMsg := pbMsg.PushMsgDataToMQ{OperationID: message.OperationID, MsgData: message.MsgData, PushToUserID: pushToUserID}
	pid, offset, err := producer.SendMessage(&mqPushMsg, mqPushMsg.PushToUserID, rpcPushMsg.OperationID)
	if err != nil {
		log.Error(mqPushMsg.OperationID, "kafka send failed", "send data", message.String(), "pid", pid, "offset", offset, "err", err.Error())
	}
	return
}

// String hashes a string to a unique hashcode.
//
// crc32 returns a uint32, but for our use we need
// and non negative integer. Here we cast to an integer
// and invert it if the result is negative.
func getHashCode(s string) uint32 {
	return crc32.ChecksumIEEE([]byte(s))
}

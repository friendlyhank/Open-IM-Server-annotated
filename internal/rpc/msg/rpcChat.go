package msg

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"github.com/golang/protobuf/proto"
)

// rpcChat - rpcChat服务端服务


// MessageWriter - 消息写入
type MessageWriter interface {
	SendMessage(m proto.Message, key string, operationID string) (int32, int64, error)
}

type rpcChat struct {
	rpcPort         int // 端口
	rpcRegisterName string // rpc注册名称
	etcdSchema      string
	etcdAddr        []string
	messageWriter   MessageWriter // 专门用于写消息
}

func NewRpcChatServer(port int) *rpcChat {
	log.NewPrivateLog(constant.LogFileName)
	rc := rpcChat{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImMsgName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
	rc.messageWriter = kafka.NewKafkaProducer(config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.Ws2mschat.Topic)
	return &rc
}

func (rpc *rpcChat) Run() {

}

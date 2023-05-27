package gate

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/statistics"
	"fmt"
	"github.com/go-playground/validator/v10"
	"sync"
)

// 初始化im和rpc相关

var (
	rwLock              *sync.RWMutex // 读写锁
	validate            *validator.Validate
	ws                  WServer   // im服务
	rpcSvr              RPCServer // rpc服务
	sendMsgAllCount     uint64    // 发送消息总数
	sendMsgFailedCount  uint64    // 发送消息失败总数
	sendMsgSuccessCount uint64    // 发送消息成功总数
	userCount           uint64    // 用户连接数

	sendMsgAllCountLock sync.RWMutex // 发送统计消息的读写锁
)

func Init(rpcPort, wsPort int) {
	rwLock = new(sync.RWMutex)
	validate = validator.New()
	statistics.NewStatistics(&sendMsgAllCount, config.Config.ModuleName.LongConnSvrName, fmt.Sprintf("%d second recv to msg_gateway sendMsgCount", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	statistics.NewStatistics(&userCount, config.Config.ModuleName.LongConnSvrName, fmt.Sprintf("%d second add user conn", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	ws.onInit(wsPort)
	rpcSvr.onInit(rpcPort) // 发送消息rpc调用
}

// Run -  运行im服务
func Run() {
	go ws.run()
	go rpcSvr.run()
}

package gate

import (
	"github.com/go-playground/validator/v10"
	"sync"
)

var (
	rwLock              *sync.RWMutex // 读写锁
	validate            *validator.Validate
	ws                  WServer
	sendMsgAllCount     uint64 // 发送消息总数
	sendMsgFailedCount  uint64 // 发送消息失败总数
	sendMsgSuccessCount uint64 // 发送消息成功总数
	userCount           uint64 // 用户连接数

	sendMsgAllCountLock sync.RWMutex // 发送统计消息的读写锁
)

func Init(rpcPort, wsPort int) {
	ws.onInit(wsPort)
}

// Run -  运行im服务
func Run() {
	go ws.run()
}

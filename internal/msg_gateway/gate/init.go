package gate

import "sync"

var (
	rwLock    *sync.RWMutex // 读写锁
	ws        WServer
	userCount uint64 // 用户连接数
)

func Init(rpcPort, wsPort int) {
	ws.onInit(wsPort)
}

// Run -  运行im服务
func Run() {
	go ws.run()
}

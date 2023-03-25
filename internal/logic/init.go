package logic

var (
	rpcServer RPCServer
)

// Init - 初始化推送服务
func Init(rpcPort int) {
	rpcServer.Init(rpcPort)
}

func init() {
}

func Run() {
	go rpcServer.run()
}

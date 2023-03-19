package msg

// rpcChat - rpcChat服务端服务
type rpcChat struct {
}

func NewRpcChatServer(port int) *rpcChat {
	rc := rpcChat{}
	return &rc
}

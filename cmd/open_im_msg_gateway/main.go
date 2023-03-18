package main

import (
	"Open_IM/internal/msg_gateway/gate"
	"Open_IM/pkg/common/config"
	"flag"
)

func main() {
	defaultRpcPorts := config.Config.RpcPort.OpenImMessageGatewayPort // rpc代理端口
	defaultWsPorts := config.Config.LongConnSvr.WebsocketPort         // socket端口
	rpcPort := flag.Int("rpc_port", defaultRpcPorts[0], "rpc listening port")
	wsPort := flag.Int("ws_port", defaultWsPorts[0], "ws listening port")
	flag.Parse()
	gate.Init(*rpcPort, *wsPort) //初始化消息代理
	gate.Run()                   // 运行im
}

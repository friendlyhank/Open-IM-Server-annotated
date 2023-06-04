package main

import (
	"Open_IM/internal/rpc/friend"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImFriendPort
	rpcPort := flag.Int("port", defaultPorts[0], "get RpcFriendPort from cmd,default 12000 as port")
	flag.Parse()
	fmt.Println("start friend rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	rpcServer := friend.NewFriendServer(*rpcPort)
	rpcServer.Run()
}

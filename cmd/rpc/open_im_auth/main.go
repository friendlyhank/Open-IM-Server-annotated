package main

import (
	rpcAuth "Open_IM/internal/rpc/auth"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImAuthPort
	rpcPort := flag.Int("port", defaultPorts[0], "RpcToken default listen port 10800")
	flag.Parse()
	fmt.Println("start auth rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	rpcServer := rpcAuth.NewRpcAuthServer(*rpcPort)
	rpcServer.Run()
}

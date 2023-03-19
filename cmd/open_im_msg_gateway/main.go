package main

import (
	"Open_IM/internal/msg_gateway/gate"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"flag"
	"fmt"
	"sync"
)

func main() {
	log.NewPrivateLog(constant.LogFileName)
	defaultRpcPorts := config.Config.RpcPort.OpenImMessageGatewayPort          // rpc代理端口
	defaultWsPorts := config.Config.LongConnSvr.WebsocketPort                  // socket端口
	defaultPromePorts := config.Config.Prometheus.MessageGatewayPrometheusPort // Prometheus监控端口
	rpcPort := flag.Int("rpc_port", defaultRpcPorts[0], "rpc listening port")
	wsPort := flag.Int("ws_port", defaultWsPorts[0], "ws listening port")
	prometheusPort := flag.Int("prometheus_port", defaultPromePorts[0], "PushrometheusPort default listen port")
	flag.Parse()
	// 要启动多个服务，所以用wait去做
	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Println("start rpc/msg_gateway server, port: ", *rpcPort, *wsPort, *prometheusPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	gate.Init(*rpcPort, *wsPort) //初始化消息代理
	gate.Run()                   // 运行im
	wg.Wait()
}

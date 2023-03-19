package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var Config config

type config struct {
	// rpc相关端口
	RpcPort struct {
		OpenImMessageGatewayPort []int `yaml:"openImMessageGatewayPort"` // im代理端口
	}
	// rpc注册的服务
	RpcRegisterName struct {
		OpenImMsgName string `yaml:"openImMsgName"`
	}
	// etcd相关配置
	Etcd struct {
		EtcdSchema string   `yaml:"etcdSchema"`
		EtcdAddr   []string `yaml:"etcdAddr"`
	}
	// 日志相关配置
	Log struct {
		RemainLogLevel uint `yaml:"remainLogLevel"`
	}
	// 长连接相关配置
	LongConnSvr struct {
		WebsocketPort       []int `yaml:"openImWsPort"`        // 端口
		WebsocketMaxConnNum int   `yaml:"websocketMaxConnNum"` // 最大连接数
		WebsocketMaxMsgLen  int   `yaml:"websocketMaxMsgLen"`  // 最大读取消息
		WebsocketTimeOut    int   `yaml:"websocketTimeOut"`    // socket连接超时时间
	}
	// Prometheus 监控
	Prometheus struct {
		MessageGatewayPrometheusPort []int `yaml:"messageGatewayPrometheusPort"`
	}
}

func unmarshalConfig(config interface{}, configName string) {
	// todo hank 先特殊处理
	bytes, err := ioutil.ReadFile(fmt.Sprintf("/Users/hank/go/src/github.com/friendlyhank/Open-IM-Server-annotated/config/%s", configName))
	if err != nil {
		panic(err.Error() + configName)
	}
	if err = yaml.Unmarshal(bytes, config); err != nil {
		panic(err.Error())
	}
}

func init() {
	unmarshalConfig(&Config, "config.yaml")
}

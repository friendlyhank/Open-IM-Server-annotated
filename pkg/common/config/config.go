package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

var(
	_, b, _, _ = runtime.Caller(0)
	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../../..")
)

var Config config

type config struct {
	// rpc相关端口
	RpcPort struct {
		OpenImMessagePort        []int `yaml:"openImMessagePort"`        // im消息端口
		OpenImMessageGatewayPort []int `yaml:"openImMessageGatewayPort"` // im代理端口
	}
	// rpc注册的服务
	RpcRegisterName struct {
		OpenImMsgName string `yaml:"openImMsgName"`
	}
	// etcd相关配置
	Etcd struct {
		StorageLocation string   `yaml:"storageLocation"`
		EtcdSchema      string   `yaml:"etcdSchema"`
		EtcdAddr        []string `yaml:"etcdAddr"`
		UserName        string   `yaml:"userName"`
		Password        string   `yaml:"password"`
	}
	// 日志相关配置
	Log struct {
		StorageLocation     string `yaml:"storageLocation"`
		RotationTime        int    `yaml:"rotationTime"`
		RemainRotationCount uint   `yaml:"remainRotationCount"`
		RemainLogLevel      uint   `yaml:"remainLogLevel"`
	}
	// 长连接相关配置
	LongConnSvr struct {
		WebsocketPort       []int `yaml:"openImWsPort"`        // 端口
		WebsocketMaxConnNum int   `yaml:"websocketMaxConnNum"` // 最大连接数
		WebsocketMaxMsgLen  int   `yaml:"websocketMaxMsgLen"`  // 最大读取消息
		WebsocketTimeOut    int   `yaml:"websocketTimeOut"`    // socket连接超时时间
	}
	// kafka相关配置
	Kafka struct {
		SASLUserName string `yaml:"SASLUserName"` // 用户
		SASLPassword string `yaml:"SASLPassword"` // 密码
	}
	// Prometheus 监控
	Prometheus struct {
		MessageGatewayPrometheusPort []int `yaml:"messageGatewayPrometheusPort"`
	}
}

func unmarshalConfig(config interface{}, configName string) {
	// todo hank 先特殊处理,研究一下路径问题
	bytes, err := ioutil.ReadFile(filepath.Join(Root, "config", configName))
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

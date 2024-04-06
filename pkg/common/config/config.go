package config

/*
 * 全局配置解析
 */

// GlobalConfig - 全局配置
type GlobalConfig struct {
	Envs struct { // 服务发现组件
		Discovery string `yaml:"discovery"`
	}
	Zookeeper struct { // zookeeper配置
		Schema   string   `yaml:"schema"`
		ZkAddr   []string `yaml:"address"`
		Username string   `yaml:"username"`
		Password string   `yaml:"password"`
	} `yaml:"zookeeper"`
	Api struct { // api端口和ip
		OpenImApiPort []int  `yaml:"openImApiPort"`
		ListenIP      string `yaml:"listenIP"`
	} `yaml:"api"`
	RpcRegisterName struct { // rpc注册服务
		OpenImUserName           string `yaml:"openImUserName"`
		OpenImFriendName         string `yaml:"openImFriendName"`
		OpenImMsgName            string `yaml:"openImMsgName"`
		OpenImPushName           string `yaml:"openImPushName"`
		OpenImMessageGatewayName string `yaml:"openImMessageGatewayName"`
		OpenImGroupName          string `yaml:"openImGroupName"`
		OpenImAuthName           string `yaml:"openImAuthName"`
		OpenImConversationName   string `yaml:"openImConversationName"`
		OpenImThirdName          string `yaml:"openImThirdName"`
	} `yaml:"rpcRegisterName"`
	Log struct { // 日志相关配置
		StorageLocation     string `yaml:"storageLocation"`
		RotationTime        uint   `yaml:"rotationTime"`
		RemainRotationCount uint   `yaml:"remainRotationCount"`
		RemainLogLevel      int    `yaml:"remainLogLevel"`
		IsStdout            bool   `yaml:"isStdout"`
		IsJson              bool   `yaml:"isJson"`
		WithStack           bool   `yaml:"withStack"`
	} `yaml:"log"`
	Prometheus struct {
		Enable            bool  `yaml:"enable"`
		ApiPrometheusPort []int `yaml:"apiPrometheusPort"`
	} `yaml:"prometheus"`
	Notification notification `yaml:"notification"` // 通知相关配置
}

type notification struct {
}

func NewGlobalConfig() *GlobalConfig {
	return &GlobalConfig{}
}

// GetServiceNames - 获取所有的rpc服务
func (c *GlobalConfig) GetServiceNames() []string {
	return []string{
		c.RpcRegisterName.OpenImUserName,
		c.RpcRegisterName.OpenImFriendName,
		c.RpcRegisterName.OpenImMsgName,
		c.RpcRegisterName.OpenImPushName,
		c.RpcRegisterName.OpenImMessageGatewayName,
		c.RpcRegisterName.OpenImGroupName,
		c.RpcRegisterName.OpenImAuthName,
		c.RpcRegisterName.OpenImConversationName,
		c.RpcRegisterName.OpenImThirdName,
	}
}

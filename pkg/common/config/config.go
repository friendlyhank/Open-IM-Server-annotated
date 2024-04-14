package config

/*
 * 全局配置解析
 */

type MYSQL struct {
	Address     []string `yaml:"address"`
	Username    string   `yaml:"username"`
	Password    string   `yaml:"password"`
	Database    string   `yaml:"database"`
	MaxOpenConn int      `yaml:"maxOpenConn"`
	MaxIdleConn int      `yaml:"maxIdleConn"`
}

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

	Mysql *MYSQL `yaml:"mysql"`

	Mongo struct {
		Uri         string   `yaml:"uri"`
		Address     []string `yaml:"address"`
		Database    string   `yaml:"database"`
		Username    string   `yaml:"username"`
		Password    string   `yaml:"password"`
		MaxPoolSize int      `yaml:"maxPoolSize"`
	} `yaml:"mongo"`

	Redis struct {
		ClusterMode    bool     `yaml:"clusterMode"` // 是否集群模式
		Address        []string `yaml:"address"`
		Username       string   `yaml:"username"`
		Password       string   `yaml:"password"`
		EnablePipeline bool     `yaml:"enablePipeline"` // 是否允许pipeline管道
	} `yaml:"redis"`
	Rpc struct { // rpc配置
		RegisterIP string `yaml:"registerIP"`
		ListenIP   string `yaml:"listenIP"`
	} `yaml:"rpc"`
	Api struct { // api端口和ip
		OpenImApiPort []int  `yaml:"openImApiPort"`
		ListenIP      string `yaml:"listenIP"`
	} `yaml:"api"`
	RpcPort struct { //rpc端口
		OpenImUserPort           []int `yaml:"openImUserPort"`
		OpenImFriendPort         []int `yaml:"openImFriendPort"`
		OpenImMessagePort        []int `yaml:"openImMessagePort"`
		OpenImMessageGatewayPort []int `yaml:"openImMessageGatewayPort"`
		OpenImGroupPort          []int `yaml:"openImGroupPort"`
		OpenImAuthPort           []int `yaml:"openImAuthPort"`
		OpenImPushPort           []int `yaml:"openImPushPort"`
		OpenImConversationPort   []int `yaml:"openImConversationPort"`
		OpenImRtcPort            []int `yaml:"openImRtcPort"`
		OpenImThirdPort          []int `yaml:"openImThirdPort"`
	} `yaml:"rpcPort"`
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
	Manager struct { // 管理员配置
		UserID   []string `yaml:"userID"`
		Nickname []string `yaml:"nickname"`
	} `yaml:"manager"`

	IMAdmin struct {
		UserID   []string `yaml:"userID"`
		Nickname []string `yaml:"nickname"`
	} `yaml:"im-admin"`

	Secret     string `yaml:"secret"`
	Prometheus struct {
		Enable                        bool   `yaml:"enable"`
		GrafanaUrl                    string `yaml:"grafanaUrl"` // grafana地址
		ApiPrometheusPort             []int  `yaml:"apiPrometheusPort"`
		UserPrometheusPort            []int  `yaml:"userPrometheusPort"`
		FriendPrometheusPort          []int  `yaml:"friendPrometheusPort"`
		MessagePrometheusPort         []int  `yaml:"messagePrometheusPort"`
		MessageGatewayPrometheusPort  []int  `yaml:"messageGatewayPrometheusPort"`
		GroupPrometheusPort           []int  `yaml:"groupPrometheusPort"`
		AuthPrometheusPort            []int  `yaml:"authPrometheusPort"`
		PushPrometheusPort            []int  `yaml:"pushPrometheusPort"`
		ConversationPrometheusPort    []int  `yaml:"conversationPrometheusPort"`
		RtcPrometheusPort             []int  `yaml:"rtcPrometheusPort"`
		MessageTransferPrometheusPort []int  `yaml:"messageTransferPrometheusPort"`
		ThirdPrometheusPort           []int  `yaml:"thirdPrometheusPort"`
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

package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../../..")
)

var Config config

type config struct {
	RpcRegisterIP string `yaml:"rpcRegisterIP"` // rpc注册ip
	ListenIP      string `yaml:"listenIP"`      // 各个rpc服务监听的ip
	// mysql 相关
	Mysql struct {
		DBAddress      []string `yaml:"dbMysqlAddress"`      // 地址
		DBUserName     string   `yaml:"dbMysqlUserName"`     // 用户名
		DBPassword     string   `yaml:"dbMysqlPassword"`     // 密码
		DBDatabaseName string   `yaml:"dbMysqlDatabaseName"` // 数据库名
		DBMaxOpenConns int      `yaml:"dbMaxOpenConns"`      // 最大开启连接
		DBMaxIdleConns int      `yaml:"dbMaxIdleConns"`      // 最大空闲连接
		DBMaxLifeTime  int      `yaml:"dbMaxLifeTime"`       // 连接存活时间
		LogLevel       int      `yaml:"logLevel"`            // 日志等级
	}
	// rpc相关端口
	RpcPort struct {
		OpenImMessagePort        []int `yaml:"openImMessagePort"`        // im消息端口
		OpenImMessageGatewayPort []int `yaml:"openImMessageGatewayPort"` // im网关端口
		OpenImPushPort           []int `yaml:"openImPushPort"`           // im消息推送端口
	}
	// rpc注册的服务
	RpcRegisterName struct {
		OpenImMsgName   string `yaml:"openImMsgName"`   // im消息名称
		OpenImPushName  string `yaml:"openImPushName"`  // 推送名称
		OpenImRelayName string `yaml:"openImRelayName"` // 消息转发,真正发送数据服务
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
		RemainLogLevel      uint   `yaml:"remainLogLevel"` // 日志等级
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
		// 消息发送topic
		Ws2mschat struct {
			Addr  []string `yaml:"addr"`  // 地址
			Topic string   `yaml:"topic"` // 对应topic
		}
		// 消费Ws2mschat的时候会触发推送，如果失败，则重新加入到推送队列
		Ms2pschat struct {
			Addr  []string `yaml:"addr"`
			Topic string   `yaml:"topic"`
		}
		// 消费者组
		ConsumerGroupID struct {
			MsgToRedis string `yaml:"msgToTransfer"`
			MsgToMongo string `yaml:"msgToMongo"` // 持久化消息到mongo
			MsgToMySql string `yaml:"msgToMySql"` // 持久化消息到mysql
		}
	}
	ChatPersistenceMysql bool `yaml:"chatpersistencemysql"` // 是否将聊天消息持久化到数据库
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

package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

/*
 * 读取加载配置
 */

var (
	_, b, _, _ = runtime.Caller(0)
	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../../..")
)

const ConfName = "openIMConf" // 配置名称

var Config config

// 回调配置信息
type callBackConfig struct {
	Enable                 bool `yaml:"enable"` // 回调开关
	CallbackTimeOut        int  `yaml:"callbackTimeOut"`
	CallbackFailedContinue bool `yaml:"callbackFailedContinue"`
}

type config struct {
	RpcRegisterIP string `yaml:"rpcRegisterIP"` // rpc注册ip
	ListenIP      string `yaml:"listenIP"`      // 各个rpc服务监听的ip
	// api服务相关
	Api struct {
		GinPort  []int  `yaml:"openImApiPort"` // gin端口设置
		ListenIP string `yaml:"listenIP"`      // ip
	}
	// cmdapi相关服务
	CmsApi struct {
		GinPort  []int  `yaml:"openImCmsApiPort"`
		ListenIP string `yaml:"listenIP"`
	}
	MultiLoginPolicy int `yaml:"multiloginpolicy"` // 多端登录配置
	TokenPolicy      struct {
		AccessSecret string `yaml:"accessSecret"` // 生成token的密钥
		AccessExpire int64  `yaml:"accessExpire"` // token过期时间
	}
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
	// redis相关配置
	Redis struct {
		DBAddress     []string `yaml:"dbAddress"`
		DBUserName    string   `yaml:"dbUserName"`
		DBPassWord    string   `yaml:"dbPassWord"`
		EnableCluster bool     `yaml:"enableCluster"`
	}
	// rpc相关端口
	RpcPort struct {
		OpenImUserPort           []int `yaml:"openImUserPort"`           // 用户相关rpc端口
		OpenImMessagePort        []int `yaml:"openImMessagePort"`        // im消息端口
		OpenImMessageGatewayPort []int `yaml:"openImMessageGatewayPort"` // im网关端口
		OpenImPushPort           []int `yaml:"openImPushPort"`           // im消息推送端口
		OpenImAuthPort           []int `yaml:"openImAuthPort"`           // rpc鉴权端口
	}
	// rpc注册的服务
	RpcRegisterName struct {
		OpenImUserName  string `yaml:"openImUserName"`  // 用户信息获取
		OpenImMsgName   string `yaml:"openImMsgName"`   // im消息名称
		OpenImPushName  string `yaml:"openImPushName"`  // 推送名称
		OpenImRelayName string `yaml:"openImRelayName"` // 消息转发,真正发送数据服务
		OpenImAuthName  string `yaml:"openImAuthName"`  // 授权验证，包含注册，获取用户token等
	}
	// etcd相关配置
	Etcd struct {
		StorageLocation string   `yaml:"storageLocation"`
		EtcdSchema      string   `yaml:"etcdSchema"`
		EtcdAddr        []string `yaml:"etcdAddr"`
		UserName        string   `yaml:"userName"`
		Password        string   `yaml:"password"`
		Secret          string   `yaml:"secret"` // 加密密钥配置
	}
	// 日志相关配置
	Log struct {
		StorageLocation     string `yaml:"storageLocation"`
		RotationTime        int    `yaml:"rotationTime"`
		RemainRotationCount uint   `yaml:"remainRotationCount"`
		RemainLogLevel      uint   `yaml:"remainLogLevel"` // 日志等级
	}
	// 模块名称
	ModuleName struct {
		LongConnSvrName string `yaml:"longConnSvrName"`
		MsgTransferName string `yaml:"msgTransferName"`
		PushName        string `yaml:"pushName"`
	}
	// 长连接相关配置
	LongConnSvr struct {
		WebsocketPort       []int `yaml:"openImWsPort"`        // 端口
		WebsocketMaxConnNum int   `yaml:"websocketMaxConnNum"` // 最大连接数
		WebsocketMaxMsgLen  int   `yaml:"websocketMaxMsgLen"`  // 最大读取消息
		WebsocketTimeOut    int   `yaml:"websocketTimeOut"`    // socket连接超时时间
	}
	ChatPersistenceMysql bool `yaml:"chatpersistencemysql"` // 是否将聊天消息持久化到数据库
	// 回调消息配置
	Callback struct {
		CallbackUrl         string         `yaml:"callbackUrl"`
		CallbackUserOnline  callBackConfig `yaml:"callbackUserOnline"`  // 用户在线回调
		CallbackUserOffline callBackConfig `yaml:"callbackUserOffline"` // 用户离线回调
		CallbackUserKickOff callBackConfig `yaml:"callbackUserKickOff"` // 用户下线回调
	}
	Demo struct {
		Port []int `yaml:"openImDemoPort"`
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
			MsgToPush  string `yaml:"msgToPush"`  // 推送消息
		}
	}
	// Prometheus 监控
	Prometheus struct {
		MessageGatewayPrometheusPort []int `yaml:"messageGatewayPrometheusPort"`
	}
}

func unmarshalConfig(config interface{}, configName string) {
	var env string
	if configName == "config.yaml" {
		env = "CONFIG_NAME"
	} else if configName == "usualConfig.yaml" {
		env = "USUAL_CONFIG_NAME"
	}
	cfgName := os.Getenv(env)
	if len(cfgName) != 0 {
		bytes, err := ioutil.ReadFile(filepath.Join(cfgName, "config", configName))
		if err != nil {
			bytes, err = ioutil.ReadFile(filepath.Join(Root, "config", configName))
			if err != nil {
				panic(err.Error() + " config: " + filepath.Join(cfgName, "config", configName))
			}
		} else {
			Root = cfgName
		}
		if err = yaml.Unmarshal(bytes, config); err != nil {
			panic(err.Error())
		}
	} else {
		// 配置是已编译的路径为启始
		bytes, err := ioutil.ReadFile(fmt.Sprintf("../config/%s", configName))
		if err != nil {
			panic(err.Error() + configName)
		}
		if err = yaml.Unmarshal(bytes, config); err != nil {
			panic(err.Error())
		}
	}
}

func init() {
	unmarshalConfig(&Config, "config.yaml")
}

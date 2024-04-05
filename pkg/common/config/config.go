package config

/*
 * 全局配置解析
 */

// GlobalConfig - 全局配置
type GlobalConfig struct {
	Api struct { // api端口和ip
		OpenImApiPort []int  `yaml:"openImApiPort"`
		ListenIP      string `yaml:"listenIP"`
	} `yaml:"api"`
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

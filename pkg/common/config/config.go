package config

/*
 * 全局配置解析
 */

// GlobalConfig - 全局配置
type GlobalConfig struct {
	Notification notification `yaml:"notification"` // 通知相关配置
}

type notification struct {
}

func NewGlobalConfig() *GlobalConfig {
	return &GlobalConfig{}
}

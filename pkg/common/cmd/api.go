package cmd

import "github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/config"

// ApiCmd - api服务程序指令
type ApiCmd struct {
	*RootCmd                                                                 //
	initFunc func(config *config.GlobalConfig, port int, promPort int) error // 初始化启动服务方法
}

func NewApiCmd() *ApiCmd {
	ret := &ApiCmd{RootCmd: NewRootCmd("api")}
	return ret
}

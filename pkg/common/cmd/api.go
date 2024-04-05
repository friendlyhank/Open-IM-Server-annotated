package cmd

import (
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/friendlyhank/open-im-server-annotated/v3/internal/api"
	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/config"
	"github.com/spf13/cobra"
)

// ApiCmd - api服务程序指令
type ApiCmd struct {
	*RootCmd                                                                 //
	initFunc func(config *config.GlobalConfig, port int, promPort int) error // 初始化启动服务方法
}

func NewApiCmd() *ApiCmd {
	ret := &ApiCmd{RootCmd: NewRootCmd("api"), initFunc: api.Start}
	ret.SetRootCmdPt(ret)
	ret.addPreRun()
	ret.addRunE()
	return ret
}

func (a *ApiCmd) addPreRun() {
	a.Command.PreRun = func(cmd *cobra.Command, args []string) {
		a.port = a.getPortFlag(cmd)
		a.prometheusPort = a.getPrometheusPortFlag(cmd)
	}
}

// addRunE - 启动程序附带错误
func (a *ApiCmd) addRunE() {
	a.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return a.initFunc(a.config, a.port, a.prometheusPort)
	}
}

func (a *ApiCmd) GetPortFromConfig(portType string) int {
	if portType == constant.FlagPort {
		return a.config.Api.OpenImApiPort[0]
	} else if portType == constant.FlagPrometheusPort {
		return a.config.Prometheus.ApiPrometheusPort[0]
	}
	return 0
}

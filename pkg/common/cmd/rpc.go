package cmd

import (
	"errors"

	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/startrpc"

	"github.com/OpenIMSDK/protocol/constant"

	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/errs"
	config2 "github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/config"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type rpcInitFuc func(config *config2.GlobalConfig, disCov discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error

type RpcCmd struct {
	*RootCmd
	RpcRegisterName string
	initFunc        rpcInitFuc
}

// NewRpcCmd - 初始化rpc指令
func NewRpcCmd(name string, initFunc rpcInitFuc) *RpcCmd {
	ret := &RpcCmd{RootCmd: NewRootCmd(name), initFunc: initFunc}
	ret.addPreRun()
	ret.addRunE()
	ret.SetRootCmdPt(ret)
	return ret
}

func (a *RpcCmd) addPreRun() {
	a.Command.PreRun = func(cmd *cobra.Command, args []string) {
		a.port = a.getPortFlag(cmd)
		a.prometheusPort = a.getPrometheusPortFlag(cmd)
	}
}

func (a *RpcCmd) addRunE() {
	a.Command.RunE = func(cmd *cobra.Command, args []string) error {
		rpcRegisterName, err := a.GetRpcRegisterNameFromConfig()
		if err != nil {
			return err
		} else {
			return a.StartSvr(rpcRegisterName, a.initFunc)
		}
	}
}

func (a *RpcCmd) Exec() error {
	return a.Execute()
}

// StartSvr - 启动rpc服务，包含服务注册与发现
func (a *RpcCmd) StartSvr(name string, rpcFn func(config *config2.GlobalConfig, disCov discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error) error {
	return startrpc.Start(a.GetPortFlag(), name, a.GetPrometheusPortFlag(), a.config, rpcFn)
}

// GetPortFromConfig - 从配置中获取端口
func (a *RpcCmd) GetPortFromConfig(portType string) int {
	switch a.Name {
	case RpcPushServer:
		if portType == constant.FlagPort {
			return a.config.RpcPort.OpenImPushPort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return a.config.Prometheus.PushPrometheusPort[0]
		}
	case RpcAuthServer:
		if portType == constant.FlagPort {
			return a.config.RpcPort.OpenImAuthPort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return a.config.Prometheus.AuthPrometheusPort[0]
		}
	case RpcConversationServer:
		if portType == constant.FlagPort {
			return a.config.RpcPort.OpenImConversationPort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return a.config.Prometheus.ConversationPrometheusPort[0]
		}
	case RpcFriendServer:
		if portType == constant.FlagPort {
			return a.config.RpcPort.OpenImFriendPort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return a.config.Prometheus.FriendPrometheusPort[0]
		}
	case RpcGroupServer:
		if portType == constant.FlagPort {
			return a.config.RpcPort.OpenImGroupPort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return a.config.Prometheus.GroupPrometheusPort[0]
		}
	case RpcMsgServer:
		if portType == constant.FlagPort {
			return a.config.RpcPort.OpenImMessagePort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return a.config.Prometheus.MessagePrometheusPort[0]
		}
	case RpcThirdServer:
		if portType == constant.FlagPort {
			return a.config.RpcPort.OpenImThirdPort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return a.config.Prometheus.ThirdPrometheusPort[0]
		}
	case RpcUserServer:
		if portType == constant.FlagPort {
			return a.config.RpcPort.OpenImUserPort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return a.config.Prometheus.UserPrometheusPort[0]
		}
	}
	return 0
}

// GetRpcRegisterNameFromConfig -从配置中获取rpc注册名
func (a *RpcCmd) GetRpcRegisterNameFromConfig() (string, error) {
	switch a.Name {
	case RpcPushServer:
		return a.config.RpcRegisterName.OpenImPushName, nil
	case RpcAuthServer:
		return a.config.RpcRegisterName.OpenImAuthName, nil
	case RpcConversationServer:
		return a.config.RpcRegisterName.OpenImConversationName, nil
	case RpcFriendServer:
		return a.config.RpcRegisterName.OpenImFriendName, nil
	case RpcGroupServer:
		return a.config.RpcRegisterName.OpenImGroupName, nil
	case RpcMsgServer:
		return a.config.RpcRegisterName.OpenImMsgName, nil
	case RpcThirdServer:
		return a.config.RpcRegisterName.OpenImThirdName, nil
	case RpcUserServer:
		return a.config.RpcRegisterName.OpenImUserName, nil
	}
	return "", errs.Wrap(errors.New("can not get rpc register name"), a.Name)
}

package main

import (
	"github.com/friendlyhank/open-im-server-annotated/v3/internal/rpc/user"
	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/cmd"
	util "github.com/friendlyhank/open-im-server-annotated/v3/pkg/util/genutil"
)

func main() {
	rpcCmd := cmd.NewRpcCmd(cmd.RpcUserServer, user.Start)
	rpcCmd.AddPortFlag()
	rpcCmd.AddPrometheusPortFlag()
	if err := rpcCmd.Exec(); err != nil {
		util.ExitWithError(err)
	}
}

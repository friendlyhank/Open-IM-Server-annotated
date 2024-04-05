package main

import (
	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/cmd"
	util "github.com/friendlyhank/open-im-server-annotated/v3/pkg/util/genutil"
)

func main() {
	apiCmd := cmd.NewApiCmd()
	apiCmd.AddPortFlag()
	apiCmd.AddPrometheusPortFlag()
	if err := apiCmd.Execute(); err != nil {
		util.ExitWithError(err)
	}
}

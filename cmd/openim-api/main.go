package main

import "github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/cmd"

func main() {
	apiCmd := cmd.NewApiCmd()
	if err := apiCmd.Execute(); err != nil {

	}
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type RootCmd struct {
	Command cobra.Command
	Name    string // 指令名称
}

// CmdOpts - 指令参数
type CmdOpts struct{}

func NewRootCmd(name string, opts ...func(*CmdOpts)) *RootCmd {
	rootCmd := &RootCmd{Name: name}
	cmd := cobra.Command{
		Use:   "Start openIM application",
		Short: fmt.Sprintf(`Start %s `, name),
		Long:  fmt.Sprintf(`Start %s `, name),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		},
	}
	rootCmd.Command = cmd
	return rootCmd
}

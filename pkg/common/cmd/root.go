package cmd

import (
	"fmt"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/config"

	"github.com/spf13/cobra"
)

type RootCmd struct {
	Command cobra.Command
	Name    string // 指令名称
	port    int
	config  *config.GlobalConfig // 全局配置
}

func (rc *RootCmd) Port() int {
	return rc.port
}

// CmdOpts - 指令参数
type CmdOpts struct{}

func NewRootCmd(name string, opts ...func(*CmdOpts)) *RootCmd {
	rootCmd := &RootCmd{Name: name, config: config.NewGlobalConfig()}
	cmd := cobra.Command{
		Use:   "Start openIM application",
		Short: fmt.Sprintf(`Start %s `, name),
		Long:  fmt.Sprintf(`Start %s `, name),
		// 执行程序前需要执行
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.persistentPreRun(cmd, opts...)
		},
	}
	rootCmd.Command = cmd
	rootCmd.addConfFlag()
	return rootCmd
}

// persistentPreRun - 执行程序前需要执行
func (rc *RootCmd) persistentPreRun(cmd *cobra.Command, opts ...func(*CmdOpts)) error {
	if err := rc.initializeConfiguration(cmd); err != nil {
		return fmt.Errorf("failed to get configuration from command: %w", err)
	}

	return nil
}

// initializeConfiguration - 初始化配置信息
func (rc *RootCmd) initializeConfiguration(cmd *cobra.Command) error {
	return rc.getConfFromCmdAndInit(cmd)
}

// addConfFlag - 添加配置文件路径参数
func (r *RootCmd) addConfFlag() {
	r.Command.Flags().StringP(constant.FlagConf, "c", "", "path to config file folder")
}

// AddPortFlag - 添加端口参数
func (r *RootCmd) AddPortFlag() {
	r.Command.Flags().IntP(constant.FlagPort, "p", 0, "server listen port")
}

// getConfFromCmdAndInit - 从命令行获取配置信息并初始化
func (r *RootCmd) getConfFromCmdAndInit(cmdLines *cobra.Command) error {
	configFolderPath, _ := cmdLines.Flags().GetString(constant.FlagConf)
	fmt.Println("The directory of the configuration file to start the process:", configFolderPath)
}

// Execute - 执行
func (r *RootCmd) Execute() error {
	return r.Command.Execute()
}

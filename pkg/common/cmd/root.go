package cmd

import (
	"fmt"

	"github.com/OpenIMSDK/tools/log"

	"github.com/OpenIMSDK/tools/errs"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/config"
	config2 "github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/config"

	"github.com/spf13/cobra"
)

// RootCmdPt - 这个不错，根据不同配置，获取每个服务对应的端口
type RootCmdPt interface {
	GetPortFromConfig(portType string) int
}

type RootCmd struct {
	Command        cobra.Command
	Name           string // 指令名称
	port           int
	prometheusPort int // prometheus端口
	cmdItf         RootCmdPt
	config         *config.GlobalConfig // 全局配置
}

func (rc *RootCmd) Port() int {
	return rc.port
}

// CmdOpts - 指令参数
type CmdOpts struct {
	loggerPrefixName string // 日志前缀
}

// WithCronTaskLogName - 初始化定时任务日志名称
func WithCronTaskLogName() func(*CmdOpts) {
	return func(opts *CmdOpts) {
		opts.loggerPrefixName = "openim-crontask"
	}
}

func WithLogName(logName string) func(*CmdOpts) {
	return func(opts *CmdOpts) {
		opts.loggerPrefixName = logName
	}
}

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

	// 追加参数信息
	cmdOpts := rc.applyOptions(opts...)

	if err := rc.initializeLogger(cmdOpts); err != nil {
		return errs.Wrap(err, "failed to initialize logger")
	}

	return nil
}

// initializeConfiguration - 初始化配置信息
func (rc *RootCmd) initializeConfiguration(cmd *cobra.Command) error {
	return rc.getConfFromCmdAndInit(cmd)
}

// applyOptions - 添加参数
func (rc *RootCmd) applyOptions(opts ...func(*CmdOpts)) *CmdOpts {
	cmdOpts := defaultCmdOpts()
	for _, opt := range opts {
		opt(cmdOpts)
	}

	return cmdOpts
}

// initializeLogger - 初始化日志
func (rc *RootCmd) initializeLogger(cmdOpts *CmdOpts) error {
	logConfig := rc.config.Log

	return log.InitFromConfig(

		cmdOpts.loggerPrefixName,
		rc.Name,
		logConfig.RemainLogLevel,
		logConfig.IsStdout,
		logConfig.IsJson,
		logConfig.StorageLocation,
		logConfig.RemainRotationCount,
		logConfig.RotationTime,
	)
}

func defaultCmdOpts() *CmdOpts {
	return &CmdOpts{
		loggerPrefixName: "openim-all",
	}
}

// SetRootCmdPt - 设置RootCmdPt
func (r *RootCmd) SetRootCmdPt(cmdItf RootCmdPt) {
	r.cmdItf = cmdItf
}

// addConfFlag - 添加配置文件路径参数
func (r *RootCmd) addConfFlag() {
	r.Command.Flags().StringP(constant.FlagConf, "c", "", "path to config file folder")
}

// AddPortFlag - 添加端口参数
func (r *RootCmd) AddPortFlag() {
	r.Command.Flags().IntP(constant.FlagPort, "p", 0, "server listen port")
}

// getPortFlag - 获取端口参数(会根据服务从配置获取)
func (r *RootCmd) getPortFlag(cmd *cobra.Command) int {
	port, err := cmd.Flags().GetInt(constant.FlagPort)
	if err != nil {
		// Wrapping the error with additional context
		return 0
	}
	if port == 0 {
		port = r.PortFromConfig(constant.FlagPort)
	}
	return port
}

// // GetPortFlag returns the port flag.
func (r *RootCmd) GetPortFlag() int {
	return r.port
}

func (r *RootCmd) AddPrometheusPortFlag() {
	r.Command.Flags().IntP(constant.FlagPrometheusPort, "", 0, "server prometheus listen port")
}

func (r *RootCmd) getPrometheusPortFlag(cmd *cobra.Command) int {
	port, err := cmd.Flags().GetInt(constant.FlagPrometheusPort)
	if err != nil || port == 0 {
		port = r.PortFromConfig(constant.FlagPrometheusPort)
		if err != nil {
			return 0
		}
	}
	return port
}

func (r *RootCmd) GetPrometheusPortFlag() int {
	return r.prometheusPort
}

// getConfFromCmdAndInit - 从命令行获取配置信息并初始化
func (r *RootCmd) getConfFromCmdAndInit(cmdLines *cobra.Command) error {
	configFolderPath, _ := cmdLines.Flags().GetString(constant.FlagConf)
	fmt.Println("The directory of the configuration file to start the process:", configFolderPath)
	return config2.InitConfig(r.config, configFolderPath)
}

// Execute - 执行
func (r *RootCmd) Execute() error {
	return r.Command.Execute()
}

func (r *RootCmd) AddCommand(cmds ...*cobra.Command) {
	r.Command.AddCommand(cmds...)
}

func (r *RootCmd) PortFromConfig(portType string) int {
	// Retrieve the port and cache it
	port := r.cmdItf.GetPortFromConfig(portType)
	return port
}

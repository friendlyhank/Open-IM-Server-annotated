package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/util/genutil"
	"gopkg.in/yaml.v3"
)

//go:embed version
var Version string

const (
	FileName             = "config.yaml"       // 配置
	NotificationFileName = "notification.yaml" // 通知配置
	DefaultFolderPath    = "../config/"
)

func GetDefaultConfigPath() string {
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Println("GetDefaultConfigPath error:", err.Error())
		return ""
	}

	configPath, err := genutil.OutDir(filepath.Join(filepath.Dir(executablePath), "../config/"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get output directory: %v\n", err)
		os.Exit(1)
	}
	return configPath
}

// GetProjectRoot - 获取项目根目录
func GetProjectRoot() string {
	executablePath, _ := os.Executable()

	projectRoot, err := genutil.OutDir(filepath.Join(filepath.Dir(executablePath), "../../../../.."))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get output directory: %v\n", err)
		os.Exit(1)
	}
	return projectRoot
}

// initConfig - 初始化配置信息
// initConfig loads configuration from a specified path into the provided config structure.
// If the specified config file does not exist, it attempts to load from the project's default "config" directory.
// It logs informative messages regarding the configuration path being used.
func initConfig(config any, configName, configFolderPath string) error {
	configFolderPath = filepath.Join(configFolderPath, configName)
	_, err := os.Stat(configFolderPath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("stat config path error:", err.Error())
			return fmt.Errorf("stat config path error: %w", err)
		}
		configFolderPath = filepath.Join(GetProjectRoot(), "config", configName)
		fmt.Println("flag's path,enviment's path,default path all is not exist,using project path:", configFolderPath)
	}
	data, err := os.ReadFile(configFolderPath)
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}
	if err = yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("unmarshal yaml error: %w", err)
	}
	fmt.Println("The path of the configuration file to start the process:", configFolderPath)

	return nil
}

// InitConfig - 初始化配置信息
func InitConfig(config *GlobalConfig, configFolderPath string) error {
	if configFolderPath == "" {
		envConfigPath := os.Getenv("OPENIMCONFIG")
		if envConfigPath != "" {
			configFolderPath = envConfigPath
		} else {
			configFolderPath = GetDefaultConfigPath()
		}
	}

	if err := initConfig(config, FileName, configFolderPath); err != nil {
		return err
	}

	return initConfig(&config.Notification, NotificationFileName, configFolderPath)
}

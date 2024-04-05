package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/util/genutil"
)

const (
	FileName             = "config.yaml"
	NotificationFileName = "notification.yaml"
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

func initConfig(config any, configName, configFolderPath string) error {
	configFolderPath = filepath.Join(configFolderPath, configName)
	_, err := os.Stat(configFolderPath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("stat config path error:", err.Error())
			return fmt.Errorf("stat config path error: %w", err)
		}
	}

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

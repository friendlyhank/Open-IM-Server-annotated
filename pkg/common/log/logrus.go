package log

import (
	"Open_IM/pkg/common/config"
	"github.com/sirupsen/logrus"
	"os"
)

var logger *Logger

type Logger struct {
	*logrus.Logger
	Pid int
}

func init() {
	logger = loggerInit("")

}

// NewPrivateLog - 初始化日志
func NewPrivateLog(moduleName string) {
	logger = loggerInit(moduleName)
}

func loggerInit(moduleName string) *Logger {
	var logger = logrus.New()
	// 设置日志级别
	logger.SetLevel(logrus.Level(config.Config.Log.RemainLogLevel))
	return &Logger{
		logger,
		os.Getpid(),
	}
}

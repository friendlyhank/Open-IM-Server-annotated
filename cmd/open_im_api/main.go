package main

import (
	apiAuth "Open_IM/internal/api/auth"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"strconv"
)

/*
 * api接口
 */

// @title open-IM-Server API
// @version 1.0
// @description  open-IM-Server 的API服务器文档, 文档中所有请求都有一个operationID字段用于链路追踪

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
func main() {
	log.NewPrivateLog(constant.LogFileName)
	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Create("../logs/api.log")
	gin.DefaultWriter = io.MultiWriter(f)
	r := gin.New()
	r.Use(gin.Recovery()) // 捕获panic
	r.Use(utils.CorsHandler())
	log.Info("load config: ", config.Config)

	//certificate 授权验证
	authRouterGroup := r.Group("/auth")
	{
		authRouterGroup.POST("/user_register", apiAuth.UserRegister) // 用户注册接口
		authRouterGroup.POST("/user_token", apiAuth.UserToken)       // 用户登录
		authRouterGroup.POST("/parse_token", apiAuth.ParseToken)     // // 解析token
	}

	defaultPorts := config.Config.Api.GinPort
	ginPort := flag.Int("port", defaultPorts[0], "get ginServerPort from cmd,default 10002 as port")
	flag.Parse()
	address := "0.0.0.0:" + strconv.Itoa(*ginPort)
	if config.Config.Api.ListenIP != "" {
		address = config.Config.Api.ListenIP + ":" + strconv.Itoa(*ginPort)
	}
	fmt.Println("start api server, address: ", address, ", OpenIM version: ", constant.CurrentVersion)
	err := r.Run(address)
	if err != nil {
		log.Error("", "api run failed ", address, err.Error())
		panic("api start failed " + err.Error())
	}
}

package main

import (
	"Open_IM/internal/custom"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"flag"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"strconv"
)

func main() {
	log.NewPrivateLog(constant.LogFileName)
	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Create("../logs/api.log")
	gin.DefaultWriter = io.MultiWriter(f)
	r := gin.Default()
	r.Use(utils.CorsHandler())
	loginRouterGroup := r.Group("/account")
	{
		loginRouterGroup.POST("/login", custom.Login)
	}
	defaultPorts := config.Config.Demo.Port
	ginPort := flag.Int("port", defaultPorts[0], "get ginServerPort from cmd,default 10004 as port")
	flag.Parse()
	address := "0.0.0.0:" + strconv.Itoa(*ginPort)
	if config.Config.Api.ListenIP != "" {
		address = config.Config.Api.ListenIP + ":" + strconv.Itoa(*ginPort)
	}
	address = config.Config.CmsApi.ListenIP + ":" + strconv.Itoa(*ginPort)
	err := r.Run(address)
	if err != nil {
		log.Error("", "run failed ", *ginPort, err.Error())
	}
}

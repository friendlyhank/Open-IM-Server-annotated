package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/rpcclient"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	util "github.com/friendlyhank/open-im-server-annotated/v3/pkg/util/genutil"

	"github.com/OpenIMSDK/tools/discoveryregistry"

	"github.com/OpenIMSDK/tools/mw"

	"github.com/gin-gonic/gin"

	"github.com/OpenIMSDK/tools/errs"

	"github.com/OpenIMSDK/tools/log"
	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/config"
	kdisc "github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/discoveryregister"
)

func Start(config *config.GlobalConfig, port int, proPort int) error {
	log.ZDebug(context.Background(), "configAPI1111111111111111111", config, "port", port, "javafdasfs")
	if port == 0 || proPort == 0 {
		err := "port or proPort is empty:" + strconv.Itoa(port) + "," + strconv.Itoa(proPort)
		return errs.Wrap(fmt.Errorf(err))
	}

	var err error

	var client discoveryregistry.SvcDiscoveryRegistry

	// Determine whether zk is passed according to whether it is a clustered deployment
	client, err = kdisc.NewDiscoveryRegister(config)
	if err != nil {
		return errs.Wrap(err, "register discovery err")
	}

	// 创建所有rpc服务的根节点
	if err = client.CreateRpcRootNodes(config.GetServiceNames()); err != nil {
		return errs.Wrap(err, "create rpc root nodes error")
	}

	var (
		netDone = make(chan struct{}, 1)
		netErr  error
	)
	router := newGinRouter(client, config)

	var address string
	if config.Api.ListenIP != "" {
		address = net.JoinHostPort(config.Api.ListenIP, strconv.Itoa(port))
	} else {
		address = net.JoinHostPort("0.0.0.0", strconv.Itoa(port))
	}

	server := http.Server{Addr: address, Handler: router}

	go func() {
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			netErr = errs.Wrap(err, fmt.Sprintf("api start err: %s", server.Addr))
			netDone <- struct{}{}
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	select {
	case <-sigs:
		util.SIGTERMExit()
		err := server.Shutdown(ctx)
		if err != nil {
			return errs.Wrap(err, "shutdown err")
		}
	case <-netDone:
		close(netDone)
		return netErr
	}
	return nil
}

// newGinRouter - 初始化gin路由
func newGinRouter(disCov discoveryregistry.SvcDiscoveryRegistry, config *config.GlobalConfig) *gin.Engine {
	// 服务发现与则册设置参数 grcp客户端，设置传输安全凭证，设置负载均衡策略为轮训的方式
	disCov.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), mw.CorsHandler(), mw.GinParseOperationID())
	// 初始化rpc客户端
	userRpc := rpcclient.NewUser(disCov, config)
	u := NewUserApi(*userRpc)
	authRpc := rpcclient.NewAuth(disCov, config)

	// 用户相关
	userRouterGroup := r.Group("/user")
	{
		userRouterGroup.POST("/user_register", u.UserRegister)
	}

	// 鉴权相关
	authRouterGroup := r.Group("/auth")
	{
		a := NewAuthApi(*authRpc)
		authRouterGroup.POST("/user_token", a.UserToken)
	}

	return r
}

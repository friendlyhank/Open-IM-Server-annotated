package startrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	util "github.com/friendlyhank/open-im-server-annotated/v3/pkg/util/genutil"

	"github.com/OpenIMSDK/tools/mw"
	"google.golang.org/grpc/credentials/insecure"

	kdisc "github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/discoveryregister"

	"github.com/OpenIMSDK/tools/errs"

	"github.com/OpenIMSDK/tools/network"

	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/config"
	config2 "github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/config"
	"google.golang.org/grpc"
)

// Start rpc server.
func Start(
	rpcPort int,
	rpcRegisterName string,
	prometheusPort int,
	config *config2.GlobalConfig,
	rpcFn func(config *config.GlobalConfig, client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error,
	options ...grpc.ServerOption,
) error {
	fmt.Printf("start %s server, port: %d, prometheusPort: %d, OpenIM version: %s\n",
		rpcRegisterName, rpcPort, prometheusPort, config2.Version)
	rpcTcpAddr := net.JoinHostPort(network.GetListenIP(config.Rpc.ListenIP), strconv.Itoa(rpcPort))
	listener, err := net.Listen(
		"tcp",
		rpcTcpAddr,
	)
	if err != nil {
		return errs.Wrap(err, "listen err", rpcTcpAddr)
	}

	defer listener.Close()
	client, err := kdisc.NewDiscoveryRegister(config)
	if err != nil {
		return err
	}

	defer client.Close()
	client.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))
	registerIP, err := network.GetRpcRegisterIP(config.Rpc.RegisterIP)
	if err != nil {
		return errs.Wrap(err)
	}

	options = append(options, mw.GrpcServer())

	srv := grpc.NewServer(options...)
	once := sync.Once{}
	defer func() {
		once.Do(srv.GracefulStop)
	}()

	err = rpcFn(config, client, srv)
	if err != nil {
		return err
	}
	// 注册服务
	err = client.Register(
		rpcRegisterName,
		registerIP,
		rpcPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return errs.Wrap(err)
	}

	var (
		netDone    = make(chan struct{}, 2)
		netErr     error
		httpServer *http.Server
	)

	go func() {
		err := srv.Serve(listener)
		if err != nil {
			netErr = errs.Wrap(err, "rpc start err: ", rpcTcpAddr)
			netDone <- struct{}{}
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	select {
	case <-sigs:
		util.SIGTERMExit()
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := gracefulStopWithCtx(ctx, srv.GracefulStop); err != nil {
			return err
		}
		ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		err := httpServer.Shutdown(ctx)
		if err != nil {
			return errs.Wrap(err, "shutdown err")
		}
		return nil
	case <-netDone:
		close(netDone)
		return netErr
	}
}

// gracefulStopWithCtx -- graceful stop with context.
func gracefulStopWithCtx(ctx context.Context, f func()) error {
	done := make(chan struct{}, 1)
	go func() {
		f()
		close(done)
	}()
	select {
	case <-ctx.Done():
		return errs.Wrap(errors.New("timeout, ctx graceful stop"))
	case <-done:
		return nil
	}
}

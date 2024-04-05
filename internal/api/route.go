package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/OpenIMSDK/tools/errs"

	"github.com/OpenIMSDK/tools/log"
	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/config"
)

func Start(config *config.GlobalConfig, port int, proPort int) error {
	log.ZDebug(context.Background(), "configAPI1111111111111111111", config, "port", port, "javafdasfs")
	if port == 0 || proPort == 0 {
		err := "port or proPort is empty:" + strconv.Itoa(port) + "," + strconv.Itoa(proPort)
		return errs.Wrap(fmt.Errorf(err))
	}

	var address string
	if config.Api.ListenIP != "" {
		address = net.JoinHostPort(config.Api.ListenIP, strconv.Itoa(port))
	} else {
		address = net.JoinHostPort("0.0.0.0", strconv.Itoa(port))
	}

	server := http.Server{Addr: address, Handler: nil}
	go func() {
		err = server.ListenAndServe()
	}()
	return nil
}

package rpcclient

import (
	"context"

	"google.golang.org/grpc"

	"github.com/OpenIMSDK/protocol/auth"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/config"
	util "github.com/friendlyhank/open-im-server-annotated/v3/pkg/util/genutil"
)

func NewAuth(discov discoveryregistry.SvcDiscoveryRegistry, config *config.GlobalConfig) *Auth {
	conn, err := discov.GetConn(context.Background(), config.RpcRegisterName.OpenImAuthName)
	if err != nil {
		util.ExitWithError(err)
	}
	client := auth.NewAuthClient(conn)
	return &Auth{discov: discov, conn: conn, Client: client, Config: config}
}

type Auth struct {
	conn   grpc.ClientConnInterface
	Client auth.AuthClient
	discov discoveryregistry.SvcDiscoveryRegistry
	Config *config.GlobalConfig
}

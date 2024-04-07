package rpcclient

import (
	"context"

	"google.golang.org/grpc"

	"github.com/OpenIMSDK/protocol/user"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/config"
	util "github.com/friendlyhank/open-im-server-annotated/v3/pkg/util/genutil"
)

// User - rpc 客户端初始化 represents a structure holding connection details for the User RPC client.
type User struct {
	conn   grpc.ClientConnInterface
	Client user.UserClient
	Discov discoveryregistry.SvcDiscoveryRegistry
	Config *config.GlobalConfig
}

// NewUser initializes and returns a User instance based on the provided service discovery registry.
func NewUser(discov discoveryregistry.SvcDiscoveryRegistry, config *config.GlobalConfig) *User {
	conn, err := discov.GetConn(context.Background(), config.RpcRegisterName.OpenImUserName)
	if err != nil {
		util.ExitWithError(err)
	}
	client := user.NewUserClient(conn)
	return &User{Discov: discov, Client: client, conn: conn, Config: config}
}

package rpcclient

import (
	"context"
	"strings"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/utils"

	"github.com/OpenIMSDK/protocol/sdkws"

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

// UserRpcClient represents the structure for a User RPC client.
type UserRpcClient User

// NewUserRpcClient initializes a UserRpcClient based on the provided service discovery registry.
func NewUserRpcClient(client discoveryregistry.SvcDiscoveryRegistry, config *config.GlobalConfig) UserRpcClient {
	return UserRpcClient(*NewUser(client, config))
}

func (u *UserRpcClient) GetUsersInfo(ctx context.Context, userIDs []string) ([]*sdkws.UserInfo, error) {
	if len(userIDs) == 0 {
		return []*sdkws.UserInfo{}, nil
	}
	resp, err := u.Client.GetDesignateUsers(ctx, &user.GetDesignateUsersReq{
		UserIDs: userIDs,
	})
	if err != nil {
		return nil, err
	}
	if ids := utils.Single(userIDs, utils.Slice(resp.UsersInfo, func(e *sdkws.UserInfo) string {
		return e.UserID
	})); len(ids) > 0 {
		return nil, errs.ErrUserIDNotFound.Wrap(strings.Join(ids, ","))
	}
	return resp.UsersInfo, nil
}

// GetUserInfo retrieves information for a single user based on the provided user ID.
func (u *UserRpcClient) GetUserInfo(ctx context.Context, userID string) (*sdkws.UserInfo, error) {
	users, err := u.GetUsersInfo(ctx, []string{userID})
	if err != nil {
		return nil, err
	}
	return users[0], nil
}

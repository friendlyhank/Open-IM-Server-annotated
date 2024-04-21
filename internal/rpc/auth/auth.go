package auth

import (
	"context"

	"github.com/OpenIMSDK/tools/errs"

	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/db/controller"

	pbauth "github.com/OpenIMSDK/protocol/auth"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/config"
	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/rpcclient"
	"google.golang.org/grpc"
)

type authServer struct {
	authDatabase  controller.AuthDatabase
	userRpcClient *rpcclient.UserRpcClient
	config        *config.GlobalConfig
}

func Start(config *config.GlobalConfig, client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	//rdb, err := cache.NewRedis(config)
	//if err != nil {
	//	return err
	//}
	userRpcClient := rpcclient.NewUserRpcClient(client, config)
	pbauth.RegisterAuthServer(server, &authServer{
		userRpcClient: &userRpcClient,
		config:        config,
		authDatabase:  controller.NewAuthDatabase(config),
	})
	return nil
}

func (s authServer) UserToken(ctx context.Context, req *pbauth.UserTokenReq) (*pbauth.UserTokenResp, error) {
	resp := pbauth.UserTokenResp{}
	if req.Secret != s.config.Secret {
		return nil, errs.ErrNoPermission.Wrap("secret invalid")
	}
	if _, err := s.userRpcClient.GetUserInfo(ctx, req.UserID); err != nil {
		return nil, err
	}
}

func (s authServer) GetUserToken(ctx context.Context, req *pbauth.GetUserTokenReq) (*pbauth.GetUserTokenResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s authServer) ForceLogout(ctx context.Context, req *pbauth.ForceLogoutReq) (*pbauth.ForceLogoutResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s authServer) ParseToken(ctx context.Context, req *pbauth.ParseTokenReq) (*pbauth.ParseTokenResp, error) {
	//TODO implement me
	panic("implement me")
}

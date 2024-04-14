package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/OpenIMSDK/protocol/constant"

	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/db/unrelation"

	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/db/cache"

	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/db/controller"

	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/utils"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"

	pbuser "github.com/OpenIMSDK/protocol/user"
	registry "github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/config"
	tablerelation "github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/db/table/relation"
	"google.golang.org/grpc"
)

// userServer - 用户服务
type userServer struct {
	controller.UserDatabase
	config *config.GlobalConfig // 配置信息
}

// Start - 启动程序
func Start(config *config.GlobalConfig, client registry.SvcDiscoveryRegistry, server *grpc.Server) error {
	rdb, err := cache.NewRedis(config)
	if err != nil {
		return err
	}
	mongo, err := unrelation.NewMongo(config)
	if err != nil {
		return err
	}
	users := make([]*tablerelation.UserModel, 0)
	if len(config.IMAdmin.UserID) != len(config.IMAdmin.Nickname) {
		return errs.Wrap(fmt.Errorf("the count of ImAdmin.UserID is not equal to the count of ImAdmin.Nickname"))
	}
	for k, v := range config.IMAdmin.UserID {
		users = append(users, &tablerelation.UserModel{UserID: v, Nickname: config.IMAdmin.Nickname[k], AppMangerLevel: constant.AppNotificationAdmin})
	}
	u := &userServer{
		config: config,
	}
	pbuser.RegisterUserServer(server, u)
	return nil
}

func (u userServer) GetDesignateUsers(ctx context.Context, req *pbuser.GetDesignateUsersReq) (*pbuser.GetDesignateUsersResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u userServer) UpdateUserInfo(ctx context.Context, req *pbuser.UpdateUserInfoReq) (*pbuser.UpdateUserInfoResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u userServer) UpdateUserInfoEx(ctx context.Context, req *pbuser.UpdateUserInfoExReq) (*pbuser.UpdateUserInfoExResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u userServer) SetGlobalRecvMessageOpt(ctx context.Context, req *pbuser.SetGlobalRecvMessageOptReq) (*pbuser.SetGlobalRecvMessageOptResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u userServer) GetGlobalRecvMessageOpt(ctx context.Context, req *pbuser.GetGlobalRecvMessageOptReq) (*pbuser.GetGlobalRecvMessageOptResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u userServer) AccountCheck(ctx context.Context, req *pbuser.AccountCheckReq) (*pbuser.AccountCheckResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u userServer) GetPaginationUsers(ctx context.Context, req *pbuser.GetPaginationUsersReq) (*pbuser.GetPaginationUsersResp, error) {
	//TODO implement me
	panic("implement me")
}

// UserRegister - 用户注册接口
func (s userServer) UserRegister(ctx context.Context, req *pbuser.UserRegisterReq) (resp *pbuser.UserRegisterResp, err error) {
	resp = &pbuser.UserRegisterResp{}
	if len(req.Users) == 0 {
		return nil, errs.ErrArgs.Wrap("users is empty")
	}
	if req.Secret != s.config.Secret {
		log.ZDebug(ctx, "UserRegister", s.config.Secret, req.Secret)
		return nil, errs.ErrNoPermission.Wrap("secret invalid")
	}
	if utils.DuplicateAny(req.Users, func(e *sdkws.UserInfo) string { return e.UserID }) {
		return nil, errs.ErrArgs.Wrap("userID repeated")
	}
	userIDs := make([]string, 0)
	for _, user := range req.Users {
		if user.UserID == "" {
			return nil, errs.ErrArgs.Wrap("userID is empty")
		}
		if strings.Contains(user.UserID, ":") {
			return nil, errs.ErrArgs.Wrap("userID contains ':' is invalid userID")
		}
		userIDs = append(userIDs, user.UserID)
	}
}

func (u userServer) GetAllUserID(ctx context.Context, req *pbuser.GetAllUserIDReq) (*pbuser.GetAllUserIDResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u userServer) UserRegisterCount(ctx context.Context, req *pbuser.UserRegisterCountReq) (*pbuser.UserRegisterCountResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u userServer) SubscribeOrCancelUsersStatus(ctx context.Context, req *pbuser.SubscribeOrCancelUsersStatusReq) (*pbuser.SubscribeOrCancelUsersStatusResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u userServer) GetSubscribeUsersStatus(ctx context.Context, req *pbuser.GetSubscribeUsersStatusReq) (*pbuser.GetSubscribeUsersStatusResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u userServer) GetUserStatus(ctx context.Context, req *pbuser.GetUserStatusReq) (*pbuser.GetUserStatusResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u userServer) SetUserStatus(ctx context.Context, req *pbuser.SetUserStatusReq) (*pbuser.SetUserStatusResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u userServer) ProcessUserCommandAdd(ctx context.Context, req *pbuser.ProcessUserCommandAddReq) (*pbuser.ProcessUserCommandAddResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u userServer) ProcessUserCommandUpdate(ctx context.Context, req *pbuser.ProcessUserCommandUpdateReq) (*pbuser.ProcessUserCommandUpdateResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u userServer) ProcessUserCommandDelete(ctx context.Context, req *pbuser.ProcessUserCommandDeleteReq) (*pbuser.ProcessUserCommandDeleteResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u userServer) ProcessUserCommandGet(ctx context.Context, req *pbuser.ProcessUserCommandGetReq) (*pbuser.ProcessUserCommandGetResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u userServer) ProcessUserCommandGetAll(ctx context.Context, req *pbuser.ProcessUserCommandGetAllReq) (*pbuser.ProcessUserCommandGetAllResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u userServer) AddNotificationAccount(ctx context.Context, req *pbuser.AddNotificationAccountReq) (*pbuser.AddNotificationAccountResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u userServer) UpdateNotificationAccountInfo(ctx context.Context, req *pbuser.UpdateNotificationAccountInfoReq) (*pbuser.UpdateNotificationAccountInfoResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u userServer) SearchNotificationAccount(ctx context.Context, req *pbuser.SearchNotificationAccountReq) (*pbuser.SearchNotificationAccountResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u userServer) GetNotificationAccount(ctx context.Context, req *pbuser.GetNotificationAccountReq) (*pbuser.GetNotificationAccountResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u userServer) GetGroupOnlineUser(ctx context.Context, req *pbuser.GetGroupOnlineUserReq) (*pbuser.GetGroupOnlineUserResp, error) {
	//TODO implement me
	panic("implement me")
}

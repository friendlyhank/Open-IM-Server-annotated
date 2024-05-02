package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/convert"

	"github.com/OpenIMSDK/tools/tx"

	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/db/mgo"

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

func (s userServer) UpdateUserInfo(ctx context.Context, req *pbuser.UpdateUserInfoReq) (*pbuser.UpdateUserInfoResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s userServer) UpdateUserInfoEx(ctx context.Context, req *pbuser.UpdateUserInfoExReq) (*pbuser.UpdateUserInfoExResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s userServer) SetGlobalRecvMessageOpt(ctx context.Context, req *pbuser.SetGlobalRecvMessageOptReq) (*pbuser.SetGlobalRecvMessageOptResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s userServer) GetGlobalRecvMessageOpt(ctx context.Context, req *pbuser.GetGlobalRecvMessageOptReq) (*pbuser.GetGlobalRecvMessageOptResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s userServer) AccountCheck(ctx context.Context, req *pbuser.AccountCheckReq) (*pbuser.AccountCheckResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s userServer) GetPaginationUsers(ctx context.Context, req *pbuser.GetPaginationUsersReq) (*pbuser.GetPaginationUsersResp, error) {
	//TODO implement me
	panic("implement me")
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
		users = append(users, &tablerelation.UserModel{UserID: v, Nickname: config.IMAdmin.Nickname[k]})
	}
	userDB, err := mgo.NewUserMongo(mongo.GetDatabase(config.Mongo.Database))
	if err != nil {
		return err
	}
	cache := cache.NewUserCacheRedis(rdb, userDB, cache.GetDefaultOpt())
	database := controller.NewUserDatabase(userDB, cache, tx.NewMongo(mongo.GetClient()))
	u := &userServer{
		UserDatabase: database,
		config:       config,
	}
	pbuser.RegisterUserServer(server, u)
	return nil
}

// GetDesignateUsers - 批量获取用户信息
func (s userServer) GetDesignateUsers(ctx context.Context, req *pbuser.GetDesignateUsersReq) (resp *pbuser.GetDesignateUsersResp, err error) {
	resp = &pbuser.GetDesignateUsersResp{}
	users, err := s.FindWithError(ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}
	resp.UsersInfo = convert.UsersDB2Pb(users)
	return resp, nil
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
	exist, err := s.IsExist(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errs.ErrRegisteredAlready.Wrap("userID registered already")
	}
	now := time.Now()
	users := make([]*tablerelation.UserModel, 0, len(req.Users))
	for _, user := range req.Users {
		users = append(users, &tablerelation.UserModel{
			UserID:     user.UserID,
			Nickname:   user.Nickname,
			FaceURL:    user.FaceURL,
			CreateTime: now,
		})
	}
	if err := s.Create(ctx, users); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s userServer) GetAllUserID(ctx context.Context, req *pbuser.GetAllUserIDReq) (*pbuser.GetAllUserIDResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s userServer) UserRegisterCount(ctx context.Context, req *pbuser.UserRegisterCountReq) (*pbuser.UserRegisterCountResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s userServer) SubscribeOrCancelUsersStatus(ctx context.Context, req *pbuser.SubscribeOrCancelUsersStatusReq) (*pbuser.SubscribeOrCancelUsersStatusResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s userServer) GetSubscribeUsersStatus(ctx context.Context, req *pbuser.GetSubscribeUsersStatusReq) (*pbuser.GetSubscribeUsersStatusResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s userServer) GetUserStatus(ctx context.Context, req *pbuser.GetUserStatusReq) (*pbuser.GetUserStatusResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s userServer) SetUserStatus(ctx context.Context, req *pbuser.SetUserStatusReq) (*pbuser.SetUserStatusResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s userServer) ProcessUserCommandAdd(ctx context.Context, req *pbuser.ProcessUserCommandAddReq) (*pbuser.ProcessUserCommandAddResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s userServer) ProcessUserCommandUpdate(ctx context.Context, req *pbuser.ProcessUserCommandUpdateReq) (*pbuser.ProcessUserCommandUpdateResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s userServer) ProcessUserCommandDelete(ctx context.Context, req *pbuser.ProcessUserCommandDeleteReq) (*pbuser.ProcessUserCommandDeleteResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s userServer) ProcessUserCommandGet(ctx context.Context, req *pbuser.ProcessUserCommandGetReq) (*pbuser.ProcessUserCommandGetResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s userServer) ProcessUserCommandGetAll(ctx context.Context, req *pbuser.ProcessUserCommandGetAllReq) (*pbuser.ProcessUserCommandGetAllResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s userServer) AddNotificationAccount(ctx context.Context, req *pbuser.AddNotificationAccountReq) (*pbuser.AddNotificationAccountResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s userServer) UpdateNotificationAccountInfo(ctx context.Context, req *pbuser.UpdateNotificationAccountInfoReq) (*pbuser.UpdateNotificationAccountInfoResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s userServer) SearchNotificationAccount(ctx context.Context, req *pbuser.SearchNotificationAccountReq) (*pbuser.SearchNotificationAccountResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s userServer) GetNotificationAccount(ctx context.Context, req *pbuser.GetNotificationAccountReq) (*pbuser.GetNotificationAccountResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s userServer) GetGroupOnlineUser(ctx context.Context, req *pbuser.GetGroupOnlineUserReq) (*pbuser.GetGroupOnlineUserResp, error) {
	//TODO implement me
	panic("implement me")
}

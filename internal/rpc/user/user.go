package user

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	sdkws "Open_IM/pkg/proto/sdk_ws"
	pbUser "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"net"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

/*
 * 用户信息rpc
 */

type userServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewUserServer(port int) *userServer {
	log.NewPrivateLog(constant.LogFileName)
	return &userServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImUserName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *userServer) Run() {
	log.NewInfo("0", "rpc user start...")

	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(s.rpcPort)

	//listener network
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("listening err:" + err.Error() + s.rpcRegisterName)
	}
	log.NewInfo("0", "listen network success, address ", address, listener)
	defer listener.Close()
	//grpc server
	var grpcOpts []grpc.ServerOption
	srv := grpc.NewServer(grpcOpts...)
	defer srv.GracefulStop()
	//Service registers with etcd
	pbUser.RegisterUserServer(srv, s)
	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("0", "RegisterEtcd failed ", err.Error(), s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName)
		panic(utils.Wrap(err, "register user module  rpc to etcd err"))
	}
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("0", "rpc  user success")
}

func (s *userServer) GetUserInfo(ctx context.Context, req *pbUser.GetUserInfoReq) (*pbUser.GetUserInfoResp, error) {
	log.NewInfo(req.OperationID, "GetUserInfo args ", req.String())
	var userInfoList []*sdkws.UserInfo
	if len(req.UserIDList) > 0 {
		for _, userID := range req.UserIDList {
			var userInfo sdkws.UserInfo
			user, err := rocksCache.GetUserInfoFromCache(userID)
			if err != nil {
				log.NewError(req.OperationID, "GetUserByUserID failed ", err.Error(), userID)
				continue
			}
			utils.CopyStructFields(&userInfo, user)
			userInfo.BirthStr = utils.TimeToString(user.Birth)
			userInfoList = append(userInfoList, &userInfo)
		}
	} else {
		return &pbUser.GetUserInfoResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: constant.ErrArgs.ErrMsg}}, nil
	}
	log.NewInfo(req.OperationID, "GetUserInfo rpc return ", pbUser.GetUserInfoResp{CommonResp: &pbUser.CommonResp{}, UserInfoList: userInfoList})
	return &pbUser.GetUserInfoResp{CommonResp: &pbUser.CommonResp{}, UserInfoList: userInfoList}, nil
}

func (s *userServer) BatchSetConversations(ctx context.Context, req *pbUser.BatchSetConversationsReq) (*pbUser.BatchSetConversationsResp, error) {
	return nil, nil
}

func (s *userServer) GetAllConversations(ctx context.Context, req *pbUser.GetAllConversationsReq) (*pbUser.GetAllConversationsResp, error) {
	return nil, nil
}

func (s *userServer) GetConversation(ctx context.Context, req *pbUser.GetConversationReq) (*pbUser.GetConversationResp, error) {
	return nil, nil
}

func (s *userServer) GetConversations(ctx context.Context, req *pbUser.GetConversationsReq) (*pbUser.GetConversationsResp, error) {
	return nil, nil
}

func (s *userServer) SetConversation(ctx context.Context, req *pbUser.SetConversationReq) (*pbUser.SetConversationResp, error) {
	return nil, nil
}

func (s *userServer) SetRecvMsgOpt(ctx context.Context, req *pbUser.SetRecvMsgOptReq) (*pbUser.SetRecvMsgOptResp, error) {
	return nil, nil
}

func (s *userServer) GetAllUserID(_ context.Context, req *pbUser.GetAllUserIDReq) (*pbUser.GetAllUserIDResp, error) {
	return nil, nil
}

func (s *userServer) AccountCheck(_ context.Context, req *pbUser.AccountCheckReq) (*pbUser.AccountCheckResp, error) {
	return nil, nil
}

func (s *userServer) UpdateUserInfo(ctx context.Context, req *pbUser.UpdateUserInfoReq) (*pbUser.UpdateUserInfoResp, error) {
	return nil, nil
}

func (s *userServer) SetGlobalRecvMessageOpt(ctx context.Context, req *pbUser.SetGlobalRecvMessageOptReq) (*pbUser.SetGlobalRecvMessageOptResp, error) {
	return nil, nil
}

func (s *userServer) GetUsers(ctx context.Context, req *pbUser.GetUsersReq) (*pbUser.GetUsersResp, error) {
	return nil, nil
}

func (s *userServer) AddUser(ctx context.Context, req *pbUser.AddUserReq) (*pbUser.AddUserResp, error) {
	return nil, nil
}

func (s *userServer) BlockUser(ctx context.Context, req *pbUser.BlockUserReq) (*pbUser.BlockUserResp, error) {
	return nil, nil
}

func (s *userServer) UnBlockUser(ctx context.Context, req *pbUser.UnBlockUserReq) (*pbUser.UnBlockUserResp, error) {
	return nil, nil
}

func (s *userServer) GetBlockUsers(ctx context.Context, req *pbUser.GetBlockUsersReq) (resp *pbUser.GetBlockUsersResp, err error) {
	return nil, nil
}

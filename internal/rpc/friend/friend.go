package friend

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbFriend "Open_IM/pkg/proto/friend"
	"Open_IM/pkg/utils"
	"context"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
	"time"
)

type friendServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewFriendServer(port int) *friendServer {
	log.NewPrivateLog(constant.LogFileName)
	return &friendServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImFriendName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *friendServer) Run() {
	log.NewInfo("0", "friendServer run...")

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
	log.NewInfo("0", "listen ok ", address)
	defer listener.Close()
	//grpc server
	var grpcOpts []grpc.ServerOption
	srv := grpc.NewServer(grpcOpts...)
	defer srv.GracefulStop()
	//User friend related services register to etcd
	pbFriend.RegisterFriendServer(srv, s)
	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	log.NewInfo("", "rpcRegisterIP", rpcRegisterIP)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("0", "RegisterEtcd failed ", err.Error(), s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName)
		panic(utils.Wrap(err, "register friend module  rpc to etcd err"))
	}
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error(), listener)
		return
	}
}

func (s *friendServer) AddBlacklist(ctx context.Context, req *pbFriend.AddBlacklistReq) (*pbFriend.AddBlacklistResp, error) {
	return nil, nil
}

// AddFriend - 添加好友
func (s *friendServer) AddFriend(ctx context.Context, req *pbFriend.AddFriendReq) (*pbFriend.AddFriendResp, error) {
	log.NewInfo(req.CommID.OperationID, "AddFriend args ", req.String())
	ok := token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID)
	if !ok {
		log.NewError(req.CommID.OperationID, "CheckAccess false ", req.CommID.OpUserID, req.CommID.FromUserID)
		return &pbFriend.AddFriendResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}
	callbackResp := callbackBeforeAddFriend(req)
	if callbackResp.ErrCode != 0 {
		log.NewError(req.CommID.OperationID, utils.GetSelfFuncName(), "callbackBeforeSendSingleMsg resp: ", callbackResp)
	}
	if callbackResp.ActionCode != constant.ActionAllow {
		if callbackResp.ErrCode == 0 {
			callbackResp.ErrCode = 201
		}
		log.NewDebug(req.CommID.OperationID, utils.GetSelfFuncName(), "callbackBeforeSendSingleMsg result", "end rpc and return", callbackResp)
		return &pbFriend.AddFriendResp{CommonResp: &pbFriend.CommonResp{
			ErrCode: int32(callbackResp.ErrCode),
			ErrMsg:  callbackResp.ErrMsg,
		}}, nil
	}
	var isSend = true
	// 目标用户的好友列表
	userIDList, err := rocksCache.GetFriendIDListFromCache(req.CommID.ToUserID)
	if err != nil {
		log.NewError(req.CommID.OperationID, "GetFriendIDListFromCache failed ", err.Error(), req.CommID.ToUserID)
		return &pbFriend.AddFriendResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	// 来源用户的好友列表
	userIDList2, err := rocksCache.GetFriendIDListFromCache(req.CommID.FromUserID)
	if err != nil {
		log.NewError(req.CommID.OperationID, "GetUserByUserID failed ", err.Error(), req.CommID.FromUserID)
		return &pbFriend.AddFriendResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	log.NewDebug(req.CommID.OperationID, "toUserID", userIDList, "fromUserID", userIDList2)
	// 如果相互之间已经是好友
	for _, v := range userIDList {
		if v == req.CommID.FromUserID {
			for _, v2 := range userIDList2 {
				if v2 == req.CommID.ToUserID {
					isSend = false
					break
				}
			}
			break
		}
	}

	//Cannot add non-existent users
	// 是否发送好友请求
	if isSend {
		if _, err := imdb.GetUserByUserID(req.CommID.ToUserID); err != nil {
			log.NewError(req.CommID.OperationID, "GetUserByUserID failed ", err.Error(), req.CommID.ToUserID)
			return &pbFriend.AddFriendResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
		friendRequest := db.FriendRequest{
			HandleResult: 0, ReqMsg: req.ReqMsg, CreateTime: time.Now()}
		utils.CopyStructFields(&friendRequest, req.CommID)
		// {openIM001 openIM002 0 test add friend 0001-01-01 00:00:00 +0000 UTC   0001-01-01 00:00:00 +0000 UTC }]
		log.NewDebug(req.CommID.OperationID, "UpdateFriendApplication args ", friendRequest)
		//err := imdb.InsertFriendApplication(&friendRequest)
		err := imdb.InsertFriendApplication(&friendRequest,
			map[string]interface{}{"handle_result": 0, "req_msg": friendRequest.ReqMsg, "create_time": friendRequest.CreateTime,
				"handler_user_id": "", "handle_msg": "", "handle_time": utils.UnixSecondToTime(0), "ex": ""})
		if err != nil {
			log.NewError(req.CommID.OperationID, "UpdateFriendApplication failed ", err.Error(), friendRequest)
			return &pbFriend.AddFriendResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}

		// todo hank 发送添加好友通知
	}
	//Establish a latest relationship in the friend request table

	return &pbFriend.AddFriendResp{CommonResp: &pbFriend.CommonResp{}}, nil
}

func (s *friendServer) ImportFriend(ctx context.Context, req *pbFriend.ImportFriendReq) (*pbFriend.ImportFriendResp, error) {
	return nil, nil
}

func (s *friendServer) AddFriendResponse(ctx context.Context, req *pbFriend.AddFriendResponseReq) (*pbFriend.AddFriendResponseResp, error) {
	return nil, nil
}

func (s *friendServer) DeleteFriend(ctx context.Context, req *pbFriend.DeleteFriendReq) (*pbFriend.DeleteFriendResp, error) {
	return nil, nil
}

func (s *friendServer) GetBlacklist(ctx context.Context, req *pbFriend.GetBlacklistReq) (*pbFriend.GetBlacklistResp, error) {
	return nil, nil
}

func (s *friendServer) SetFriendRemark(ctx context.Context, req *pbFriend.SetFriendRemarkReq) (*pbFriend.SetFriendRemarkResp, error) {
	return nil, nil
}

func (s *friendServer) RemoveBlacklist(ctx context.Context, req *pbFriend.RemoveBlacklistReq) (*pbFriend.RemoveBlacklistResp, error) {
	return nil, nil
}

func (s *friendServer) IsInBlackList(ctx context.Context, req *pbFriend.IsInBlackListReq) (*pbFriend.IsInBlackListResp, error) {
	return nil, nil
}

func (s *friendServer) IsFriend(ctx context.Context, req *pbFriend.IsFriendReq) (*pbFriend.IsFriendResp, error) {
	return nil, nil
}

func (s *friendServer) GetFriendList(ctx context.Context, req *pbFriend.GetFriendListReq) (*pbFriend.GetFriendListResp, error) {
	return nil, nil
}

func (s *friendServer) GetFriendApplyList(ctx context.Context, req *pbFriend.GetFriendApplyListReq) (*pbFriend.GetFriendApplyListResp, error) {
	return nil, nil
}

func (s *friendServer) GetSelfApplyList(ctx context.Context, req *pbFriend.GetSelfApplyListReq) (*pbFriend.GetSelfApplyListResp, error) {
	return nil, nil
}

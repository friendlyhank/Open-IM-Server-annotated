package auth

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbAuth "Open_IM/pkg/proto/auth"
	"Open_IM/pkg/utils"
	"context"
	"net"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

// UserRegister - 用户注册
func (rpc *rpcAuth) UserRegister(_ context.Context, req *pbAuth.UserRegisterReq) (*pbAuth.UserRegisterResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	var user db.User
	utils.CopyStructFields(&user, req.UserInfo)
	if req.UserInfo.BirthStr != "" {
		time, err := utils.TimeStringToTime(req.UserInfo.BirthStr)
		if err != nil {
			log.NewError(req.OperationID, "TimeStringToTime failed ", err.Error(), req.UserInfo.BirthStr)
			return &pbAuth.UserRegisterResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: "TimeStringToTime failed:" + err.Error()}}, nil
		}
		user.Birth = time
	}
	log.Debug(req.OperationID, "copy ", user, req.UserInfo)
	err := imdb.UserRegister(user)
	if err != nil {
		errMsg := req.OperationID + " imdb.UserRegister failed " + err.Error() + user.UserID
		log.NewError(req.OperationID, errMsg, user)
		return &pbAuth.UserRegisterResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}}, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ", pbAuth.UserRegisterResp{CommonResp: &pbAuth.CommonResp{}})
	return &pbAuth.UserRegisterResp{CommonResp: &pbAuth.CommonResp{}}, nil
}

// UserToken - 获取用户token
func (rpc *rpcAuth) UserToken(_ context.Context, req *pbAuth.UserTokenReq) (*pbAuth.UserTokenResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	_, err := imdb.GetUserByUserID(req.FromUserID)
	if err != nil {
		log.NewError(req.OperationID, "not this user:", req.FromUserID, req.String())
		return &pbAuth.UserTokenResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	tokens, expTime, err := token_verify.CreateToken(req.FromUserID, int(req.Platform))
	if err != nil {
		errMsg := req.OperationID + " token_verify.CreateToken failed " + err.Error() + req.FromUserID + utils.Int32ToString(req.Platform)
		log.NewError(req.OperationID, errMsg)
		return &pbAuth.UserTokenResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}}, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ", pbAuth.UserTokenResp{CommonResp: &pbAuth.CommonResp{}, Token: tokens, ExpiredTime: expTime})
	return &pbAuth.UserTokenResp{CommonResp: &pbAuth.CommonResp{}, Token: tokens, ExpiredTime: expTime}, nil
}

func (rpc *rpcAuth) ParseToken(_ context.Context, req *pbAuth.ParseTokenReq) (*pbAuth.ParseTokenResp, error) {
	return nil, nil
}

func (rpc *rpcAuth) ForceLogout(_ context.Context, req *pbAuth.ForceLogoutReq) (*pbAuth.ForceLogoutResp, error) {
	return nil, nil
}

type rpcAuth struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewRpcAuthServer(port int) *rpcAuth {
	log.NewPrivateLog(constant.LogFileName)
	return &rpcAuth{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImAuthName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (rpc *rpcAuth) Run() {
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, "rpc auth start...")

	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(rpc.rpcPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("listening err:" + err.Error() + rpc.rpcRegisterName)
	}
	log.NewInfo(operationID, "listen network success, ", address, listener)
	var grpcOpts []grpc.ServerOption
	srv := grpc.NewServer(grpcOpts...)
	defer srv.GracefulStop()

	//service registers with etcd
	pbAuth.RegisterAuthServer(srv, rpc)
	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	log.NewInfo("", "rpcRegisterIP", rpcRegisterIP)

	err = getcdv3.RegisterEtcd(rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName, 10)
	if err != nil {
		log.NewError(operationID, "RegisterEtcd failed ", err.Error(),
			rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName)
		panic(utils.Wrap(err, "register auth module  rpc to etcd err"))

	}
	log.NewInfo(operationID, "RegisterAuthServer ok ", rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName)
	err = srv.Serve(listener)
	if err != nil {
		log.NewError(operationID, "Serve failed ", err.Error())
		return
	}
	log.NewInfo(operationID, "rpc auth ok")
}

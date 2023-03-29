package gate

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbRelay "Open_IM/pkg/proto/relay"
	"Open_IM/pkg/utils"
	"context"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
)

// rpc服务相关逻辑
type RPCServer struct {
	rpcPort         int // rpc端口
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string // etcd地址
}

// rpc服务初始化
func (r *RPCServer) onInit(rpcPort int) {
	r.rpcPort = rpcPort
	r.etcdSchema = config.Config.Etcd.EtcdSchema
	r.etcdAddr = config.Config.Etcd.EtcdAddr
}

func (r *RPCServer) run() {
	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(r.rpcPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("listening err:" + err.Error() + r.rpcRegisterName)
	}
	defer listener.Close()
	var grpcOpts []grpc.ServerOption
	srv := grpc.NewServer(grpcOpts...)
	defer srv.GracefulStop()
	pbRelay.RegisterRelayServer(srv, r)

	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	err = getcdv3.RegisterEtcd4Unique(r.etcdSchema, strings.Join(r.etcdAddr, ","), rpcRegisterIP, r.rpcPort, r.rpcRegisterName, 10)
	if err != nil {
		log.Error("", "register push message rpc to etcd err", "", "err", err.Error(), r.etcdSchema, strings.Join(r.etcdAddr, ","), rpcRegisterIP, r.rpcPort, r.rpcRegisterName)
		panic(utils.Wrap(err, "register msg_gataway module  rpc to etcd err"))
	}
	err = srv.Serve(listener)
	if err != nil {
		log.Error("", "push message rpc listening err", "", "err", err.Error())
		return
	}
	err = srv.Serve(listener)
	if err != nil {
		log.Error("", "push message rpc listening err", "", "err", err.Error())
		return
	}
}

func (r *RPCServer) OnlinePushMsg(_ context.Context, in *pbRelay.OnlinePushMsgReq) (*pbRelay.OnlinePushMsgResp, error) {
	return nil, nil
}

// GetUsersOnlineStatus - 获取用户在线状态
func (r *RPCServer) GetUsersOnlineStatus(_ context.Context, req *pbRelay.GetUsersOnlineStatusReq) (*pbRelay.GetUsersOnlineStatusResp, error) {
	return nil, nil
}

func (r *RPCServer) SuperGroupOnlineBatchPushOneMsg(_ context.Context, req *pbRelay.OnlineBatchPushOneMsgReq) (*pbRelay.OnlineBatchPushOneMsgResp, error) {
	return nil, nil
}

func (r *RPCServer) SuperGroupBackgroundOnlinePush(_ context.Context, req *pbRelay.OnlineBatchPushOneMsgReq) (*pbRelay.OnlineBatchPushOneMsgResp, error) {
	return nil, nil
}

func (r *RPCServer) OnlineBatchPushOneMsg(_ context.Context, req *pbRelay.OnlineBatchPushOneMsgReq) (*pbRelay.OnlineBatchPushOneMsgResp, error) {
	return nil, nil
}

// KickUserOffline - 把用户踢下线
func (r *RPCServer) KickUserOffline(_ context.Context, req *pbRelay.KickUserOfflineReq) (*pbRelay.KickUserOfflineResp, error) {
	return nil, nil
}

func (r *RPCServer) MultiTerminalLoginCheck(ctx context.Context, req *pbRelay.MultiTerminalLoginCheckReq) (*pbRelay.MultiTerminalLoginCheckResp, error) {
	return nil, nil
}

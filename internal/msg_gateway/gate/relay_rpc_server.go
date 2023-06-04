package gate

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbRelay "Open_IM/pkg/proto/relay"
	"Open_IM/pkg/utils"
	"bytes"
	"context"
	"encoding/gob"
	"net"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

// rpc服务相关逻辑
type RPCServer struct {
	rpcPort         int // rpc端口
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string // etcd地址
	pushTerminal    []int    // 推送的端(platformID)
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

// SuperGroupOnlineBatchPushOneMsg - 超级管理员批量推送消息
func (r *RPCServer) SuperGroupOnlineBatchPushOneMsg(_ context.Context, req *pbRelay.OnlineBatchPushOneMsgReq) (*pbRelay.OnlineBatchPushOneMsgResp, error) {
	log.NewInfo(req.OperationID, "BatchPushMsgToUser is arriving", req.String())
	var singleUserResult []*pbRelay.SingelMsgToUserResultList
	msgBytes, _ := proto.Marshal(req.MsgData)
	mReply := Resp{
		ReqIdentifier: constant.WSPushMsg,
		OperationID:   req.OperationID,
		Data:          msgBytes,
	}
	var replyBytes bytes.Buffer
	enc := gob.NewEncoder(&replyBytes)
	err := enc.Encode(mReply)
	if err != nil {
		log.NewError(req.OperationID, "data encode err", err.Error())
	}
	for _, v := range req.PushToUserIDList {
		var resp []*pbRelay.SingleMsgToUserPlatform
		tempT := &pbRelay.SingelMsgToUserResultList{
			UserID: v,
		}
		userConnMap := ws.getUserAllCons(v)
		for platform, userConns := range userConnMap {
			if userConns != nil {
				log.NewWarn(req.OperationID, "conns is ", len(userConns), platform, userConns)
				for _, userConn := range userConns {
					temp := &pbRelay.SingleMsgToUserPlatform{
						RecvID:         v,
						RecvPlatFormID: int32(platform),
					}
					resultCode := sendMsgBatchToUser(userConn, replyBytes.Bytes(), req, platform, v)
					if resultCode == 0 {
						tempT.OnlinePush = true
						temp.ResultCode = resultCode
						resp = append(resp, temp)
					}
				}
			}
		}
		tempT.Resp = resp
		singleUserResult = append(singleUserResult, tempT)
	}
	return &pbRelay.OnlineBatchPushOneMsgResp{
		SinglePushResult: singleUserResult,
	}, nil
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

// sendMsgBatchToUser - 发送消息给用户
func sendMsgBatchToUser(conn *UserConn, bMsg []byte, in *pbRelay.OnlineBatchPushOneMsgReq, RecvPlatForm int, RecvID string) (ResultCode int64) {
	err := ws.writeMsg(conn, websocket.BinaryMessage, bMsg)
	if err != nil {
		log.NewError(in.OperationID, "PushMsgToUser is failed By Ws", "Addr", conn.RemoteAddr().String(),
			"error", err, "senderPlatform", constant.PlatformIDToName(int(in.MsgData.SenderPlatformID)), "recv Platform", RecvPlatForm, "args", in.String(), "recvID", RecvID)
		ResultCode = -2
		return ResultCode
	} else {
		log.NewDebug(in.OperationID, "PushMsgToUser is success By Ws", "args", in.String(), "recv PlatForm", RecvPlatForm, "recvID", RecvID)
		ResultCode = 0
		return ResultCode
	}
}

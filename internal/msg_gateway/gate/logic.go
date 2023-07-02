package gate

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbChat "Open_IM/pkg/proto/msg"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"bytes"
	"context"
	"encoding/gob"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"strings"
)

// 处理im消息的逻辑

// msgParse - 解析消息
func (ws *WServer) msgParse(conn *UserConn, binaryMsg []byte) {
	b := bytes.NewBuffer(binaryMsg)
	m := Req{}
	dec := gob.NewDecoder(b)
	err := dec.Decode(&m)
	if err != nil {
		log.NewError("", "ws Decode  err", err.Error())
		err = conn.Close()
		if err != nil {
			log.NewError("", "ws close err", err.Error())
		}
		return
	}
	if err := validate.Struct(m); err != nil {
		log.NewError("", "ws args validate  err", err.Error())
		ws.sendErrMsg(conn, 201, err.Error(), m.ReqIdentifier, m.MsgIncr, m.OperationID)
		return
	}
	log.NewInfo(m.OperationID, "Basic Info Authentication Success", m.SendID, m.MsgIncr, m.ReqIdentifier)
	if m.SendID != conn.userID {
		if err = conn.Close(); err != nil {
			log.NewError(m.OperationID, "close ws conn failed", conn.userID, "send id", m.SendID, err.Error())
			return
		}
	}
	switch m.ReqIdentifier {
	case constant.WSGetNewestSeq:
		log.NewInfo(m.OperationID, "getSeqReq ", m.SendID, m.MsgIncr, m.ReqIdentifier)
		ws.getSeqReq(conn, &m)
	case constant.WSSendMsg: // 发送消息
		log.NewInfo(m.OperationID, "sendMsgReq ", m.SendID, m.MsgIncr, m.ReqIdentifier)
		ws.sendMsgReq(conn, &m)
	}
}

// getSeqReq - 获取最新的req序号
func (ws *WServer) getSeqReq(conn *UserConn, m *Req) {
	log.NewInfo(m.OperationID, "Ws call success to getNewSeq", m.MsgIncr, m.SendID, m.ReqIdentifier)
	nReply := new(sdk_ws.GetMaxAndMinSeqResp)
	isPass, errCode, errMsg, data := ws.argsValidate(m, constant.WSGetNewestSeq, m.OperationID)
	log.Info(m.OperationID, "argsValidate ", isPass, errCode, errMsg)
	if isPass {
		rpcReq := sdk_ws.GetMaxAndMinSeqReq{}
		// todo hank groupid
		rpcReq.UserID = m.SendID
		rpcReq.OperationID = m.OperationID
		log.Debug(m.OperationID, "Ws call success to getMaxAndMinSeq", m.SendID, m.ReqIdentifier, m.MsgIncr, data.(sdk_ws.GetMaxAndMinSeqReq).GroupIDList)
		grpcConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, rpcReq.OperationID)
		if grpcConn == nil {
			errMsg := rpcReq.OperationID + "getcdv3.GetDefaultConn == nil"
			nReply.ErrCode = 500
			nReply.ErrMsg = errMsg
			log.NewError(rpcReq.OperationID, errMsg)
			ws.getSeqResp(conn, m, nReply)
			return
		}
		msgClient := pbChat.NewMsgClient(grpcConn)
		rpcReply, err := msgClient.GetMaxAndMinSeq(context.Background(), &rpcReq)
		if err != nil {
			nReply.ErrCode = 500
			nReply.ErrMsg = err.Error()
			log.Error(rpcReq.OperationID, "rpc call failed to GetMaxAndMinSeq ", nReply.String())
			ws.getSeqResp(conn, m, nReply)
		} else {
			log.NewInfo(rpcReq.OperationID, "rpc call success to getSeqReq", rpcReply.String())
			ws.getSeqResp(conn, m, rpcReply)
		}
	} else {
		nReply.ErrCode = errCode
		nReply.ErrMsg = errMsg
		log.Error(m.OperationID, "argsValidate failed send resp: ", nReply.String())
		ws.getSeqResp(conn, m, nReply)
	}
}

// getSeqResp - 获取最大序号req答复
func (ws *WServer) getSeqResp(conn *UserConn, m *Req, pb *sdk_ws.GetMaxAndMinSeqResp) {

}

// sendMsgReq - 发送消息请求
func (ws *WServer) sendMsgReq(conn *UserConn, m *Req) {
	sendMsgAllCountLock.Lock()
	sendMsgAllCount++
	sendMsgAllCountLock.Unlock()
	log.NewInfo(m.OperationID, "Ws call success to sendMsgReq start", m.MsgIncr, m.ReqIdentifier, m.SendID)

	nReply := new(pbChat.SendMsgResp)
	// 参数校验
	isPass, errCode, errMsg, pData := ws.argsValidate(m, constant.WSSendMsg, m.OperationID)
	if isPass {
		data := pData.(sdk_ws.MsgData)
		pbData := pbChat.SendMsgReq{
			Token:       m.Token,
			OperationID: m.OperationID,
			MsgData:     &data,
		}
		log.NewInfo(m.OperationID, "Ws call success to sendMsgReq middle", m.ReqIdentifier, m.SendID, m.MsgIncr, data.String())
		// 获取连接
		etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, m.OperationID)
		if etcdConn == nil {
			errMsg := m.OperationID + "getcdv3.GetDefaultConn == nil"
			nReply.ErrCode = 500
			nReply.ErrMsg = errMsg
			log.NewError(m.OperationID, errMsg)
			ws.sendMsgResp(conn, m, nReply)
			return
		}
		// rpc调用服务
		client := pbChat.NewMsgClient(etcdConn)
		reply, err := client.SendMsg(context.Background(), &pbData)
		if err != nil {
			nReply.ErrCode = 200
			nReply.ErrMsg = err.Error()
			ws.sendMsgResp(conn, m, nReply)
		} else {
			log.NewInfo(pbData.OperationID, "rpc call success to sendMsgReq", reply.String())
			ws.sendMsgResp(conn, m, reply)
		}
	} else {
		nReply.ErrCode = errCode
		nReply.ErrMsg = errMsg
		ws.sendMsgResp(conn, m, nReply)
	}
}

// sendMsgResp - 发送消息响应
func (ws *WServer) sendMsgResp(conn *UserConn, m *Req, pb *pbChat.SendMsgResp) {
	var mReplyData sdk_ws.UserSendMsgResp
	mReplyData.ClientMsgID = pb.GetClientMsgID()
	mReplyData.ServerMsgID = pb.GetServerMsgID()
	mReplyData.SendTime = pb.GetSendTime()
	b, _ := proto.Marshal(&mReplyData)
	mReply := Resp{
		ReqIdentifier: m.ReqIdentifier,
		ErrCode:       pb.GetErrCode(),
		ErrMsg:        pb.GetErrMsg(),
		OperationID:   m.OperationID,
		Data:          b,
	}
	ws.sendMsg(conn, mReply)
}

// sendMsg - 发送答复消息
func (ws *WServer) sendMsg(conn *UserConn, mReply interface{}) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(mReply)
	if err != nil {
		//	uid, platform := ws.getUserUid(conn)
		log.NewError(mReply.(Resp).OperationID, mReply.(Resp).ReqIdentifier, mReply.(Resp).ErrCode, mReply.(Resp).ErrMsg, "Encode Msg error", conn.RemoteAddr().String(), err.Error())
		return
	}
	err = ws.writeMsg(conn, websocket.BinaryMessage, b.Bytes())
	if err != nil {
		//	uid, platform := ws.getUserUid(conn)
		log.NewError(mReply.(Resp).OperationID, mReply.(Resp).ReqIdentifier, mReply.(Resp).ErrCode, mReply.(Resp).ErrMsg, "ws writeMsg error", conn.RemoteAddr().String(), err.Error())
	} else {
		log.Debug(mReply.(Resp).OperationID, mReply.(Resp).ReqIdentifier, mReply.(Resp).ErrCode, mReply.(Resp).ErrMsg, "ws write response success")
	}
}

// sendErrMsg - 发送错误的消息
func (ws *WServer) sendErrMsg(conn *UserConn, errCode int32, errMsg string, reqIdentifier int32, msgIncr string, operationID string) {
	mReply := Resp{
		ReqIdentifier: reqIdentifier,
		MsgIncr:       msgIncr,
		ErrCode:       errCode,
		ErrMsg:        errMsg,
		OperationID:   operationID,
	}
	ws.sendMsg(conn, mReply)
}

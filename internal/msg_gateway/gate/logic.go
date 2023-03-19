package gate

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
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
		err = conn.Close()
		if err != nil {
		}
		return
	}
	if m.SendID != conn.userID {
		if err = conn.Close(); err != nil {
			return
		}
	}
	switch m.ReqIdentifier {
	case constant.WSSendMsg:
	}
}

// sendMsgReq - 发送消息请求
func (ws *WServer) sendMsgReq(conn *UserConn, m *Req) {
	sendMsgAllCountLock.Lock()
	sendMsgAllCount++
	sendMsgAllCountLock.Unlock()

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
		// 获取连接
		etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, m.OperationID)
		if etcdConn == nil {
			errMsg := m.OperationID + "getcdv3.GetDefaultConn == nil"
			nReply.ErrCode = 500
			nReply.ErrMsg = errMsg
			ws.sendMsgResp(conn, m, nReply)
		}
		// rpc调用服务
		client := pbChat.NewMsgClient(etcdConn)
		reply, err := client.SendMsg(context.Background(), &pbData)
		if err != nil {
			nReply.ErrCode = 200
			nReply.ErrMsg = err.Error()
			ws.sendMsgResp(conn, m, nReply)
		} else {
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
		return
	}
	err = ws.writeMsg(conn, websocket.BinaryMessage, b.Bytes())
	if err != nil {

	} else {

	}
}

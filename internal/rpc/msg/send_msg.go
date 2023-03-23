package msg

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	pbChat "Open_IM/pkg/proto/msg"
	"context"
	"errors"
	"time"
)

// 消息发送server逻辑

// SendMsg - 发送消息
func (rpc *rpcChat) SendMsg(_ context.Context, pb *pbChat.SendMsgReq) (*pbChat.SendMsgResp, error) {
	replay := pbChat.SendMsgResp{}

	t1 := time.Now()
	// 构建发送mq的消息
	msgToMQSingle := pbChat.MsgDataToMQ{Token: pb.Token, OperationID: pb.OperationID, MsgData: pb.MsgData}

	log.Info(pb.OperationID, "rpc sendMsg come here ", pb.String())
	switch pb.MsgData.SessionType {
	case constant.SingleChatType: // 单聊消息
		t1 = time.Now()
		if msgToMQSingle.MsgData.SendID != msgToMQSingle.MsgData.RecvID { //Filter messages sent to yourself
			t1 = time.Now()
		}
		err2 := rpc.sendMsgToWriter(&msgToMQSingle, msgToMQSingle.MsgData.SendID, constant.OnlineStatus)
		log.Info(pb.OperationID, "sendMsgToWriter ", " cost time: ", time.Since(t1))
		if err2 != nil {
			log.NewError(msgToMQSingle.OperationID, "kafka send msg err:SendID", msgToMQSingle.MsgData.SendID, msgToMQSingle.String())
			return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
		}

	default:
		return returnMsg(&replay, pb, 203, "unknown sessionType", "", 0)
	}
	return &replay,nil
}

// returnMsg - 返回发送消息答复
func returnMsg(replay *pbChat.SendMsgResp, pb *pbChat.SendMsgReq, errCode int32, errMsg, serverMsgID string, sendTime int64) (*pbChat.SendMsgResp, error) {
	replay.ErrCode = errCode
	replay.ErrMsg = errMsg
	replay.ServerMsgID = serverMsgID
	replay.ClientMsgID = pb.MsgData.ClientMsgID
	replay.SendTime = sendTime
	return replay, nil
}

// sendMsgToWriter - 发送消息到writer
func (rpc *rpcChat) sendMsgToWriter(m *pbChat.MsgDataToMQ, key string, status string) error {
	switch status {
	case constant.OnlineStatus:
		pid, offset, err := rpc.messageWriter.SendMessage(m, key, m.OperationID)
		if err != nil {
			log.Error(m.OperationID, "kafka send failed", "send data", m.String(), "pid", pid, "offset", offset, "err", err.Error(), "key", key, status)
		} else {
			//	log.NewWarn(m.OperationID, "sendMsgToWriter   client msgID ", m.MsgData.ClientMsgID)
		}
		return err
	case constant.OfflineStatus:
	}
	return errors.New("status error")
}
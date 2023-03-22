package msg

import (
	"Open_IM/pkg/common/log"
	pbChat "Open_IM/pkg/proto/msg"
	"context"
)

// 消息发送server逻辑

// SendMsg - 发送消息
func (rpc *rpcChat) SendMsg(_ context.Context, pb *pbChat.SendMsgReq) (*pbChat.SendMsgResp, error) {
	// replay := pbChat.SendMsgResp{}
	log.Info(pb.OperationID, "rpc sendMsg come here ", pb.String())
	return nil,nil
}

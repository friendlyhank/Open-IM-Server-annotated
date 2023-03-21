package msg

import (
	pbChat "Open_IM/pkg/proto/msg"
	"Open_IM/pkg/common/log"
	"context"
)

// 消息发送server逻辑

// SendMsg - 发送消息
func (rpc *rpcChat) SendMsg(_ context.Context, pb *pbChat.SendMsgReq) (*pbChat.SendMsgResp, error) {
	replay := pbChat.SendMsgResp{}
	log.Info(pb.OperationID, "rpc sendMsg come here ", pb.String())

	return nil,nil
}

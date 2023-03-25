package msg

import (
	pbMsg "Open_IM/pkg/proto/msg"
	"context"
)

func (rpc *rpcChat) SetSendMsgStatus(_ context.Context, req *pbMsg.SetSendMsgStatusReq) (resp *pbMsg.SetSendMsgStatusResp, err error) {
	return nil,nil
}

func (rpc *rpcChat) GetSendMsgStatus(_ context.Context, req *pbMsg.GetSendMsgStatusReq) (resp *pbMsg.GetSendMsgStatusResp, err error) {
	return nil,nil
}

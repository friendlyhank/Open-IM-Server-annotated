package msg

import (
	pbChat "Open_IM/pkg/proto/msg"
	"context"
)

func (rpc *rpcChat) ClearMsg(_ context.Context, req *pbChat.ClearMsgReq) (*pbChat.ClearMsgResp, error) {
	return nil,nil
}

func (rpc *rpcChat) SetMsgMinSeq(_ context.Context, req *pbChat.SetMsgMinSeqReq) (*pbChat.SetMsgMinSeqResp, error) {
	return nil,nil
}

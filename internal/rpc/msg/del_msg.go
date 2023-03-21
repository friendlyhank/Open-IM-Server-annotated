package msg

import (
	"Open_IM/pkg/proto/msg"
	commonPb "Open_IM/pkg/proto/sdk_ws"
	"context"
)

func (rpc *rpcChat) DelMsgList(_ context.Context, req *commonPb.DelMsgListReq) (*commonPb.DelMsgListResp, error) {
	return nil,nil
}

func (rpc *rpcChat) DelSuperGroupMsg(_ context.Context, req *msg.DelSuperGroupMsgReq) (*msg.DelSuperGroupMsgResp, error) {
	return nil,nil
}

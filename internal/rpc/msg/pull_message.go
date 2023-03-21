package msg

import (
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"context"
)

func (rpc *rpcChat) GetMaxAndMinSeq(_ context.Context, in *open_im_sdk.GetMaxAndMinSeqReq) (*open_im_sdk.GetMaxAndMinSeqResp, error) {
	return nil,nil
}

func (rpc *rpcChat) PullMessageBySeqList(_ context.Context, in *open_im_sdk.PullMessageBySeqListReq) (*open_im_sdk.PullMessageBySeqListResp, error) {
	return nil,nil
}

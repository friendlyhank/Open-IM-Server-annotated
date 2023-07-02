package msg

import (
	commonDB "Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"context"
)

func (rpc *rpcChat) GetMaxAndMinSeq(_ context.Context, in *open_im_sdk.GetMaxAndMinSeqReq) (*open_im_sdk.GetMaxAndMinSeqResp, error) {
	log.NewInfo(in.OperationID, "rpc getMaxAndMinSeq is arriving", in.String())
	resp := new(open_im_sdk.GetMaxAndMinSeqResp)
	m := make(map[string]*open_im_sdk.MaxAndMinSeq)
	// 获取最大的seq和最小的seq
	var maxSeq, minSeq uint64
	var err1, err2 error
	maxSeq, err1 = commonDB.DB.GetUserMaxSeq(in.UserID)
	minSeq, err2 = commonDB.DB.GetUserMinSeq(in.UserID)
}

func (rpc *rpcChat) PullMessageBySeqList(_ context.Context, in *open_im_sdk.PullMessageBySeqListReq) (*open_im_sdk.PullMessageBySeqListResp, error) {
	return nil, nil
}

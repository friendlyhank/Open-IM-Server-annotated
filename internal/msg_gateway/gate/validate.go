/*
** description("").
** copyright('Open_IM,www.Open_IM.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/5/21 15:29).
 */
package gate

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"github.com/golang/protobuf/proto"
)

// Req - im请求结构体
type Req struct {
	ReqIdentifier int32  `json:"reqIdentifier" validate:"required"` // 请求标识，对应socket协议
	Token         string `json:"token" `
	SendID        string `json:"sendID" validate:"required"` // 发送id
	OperationID   string `json:"operationID" validate:"required"`
	MsgIncr       string `json:"msgIncr" validate:"required"` // todo hank 未知
	Data          []byte `json:"data"`
}

// Resp - im答复结构体
type Resp struct {
	ReqIdentifier int32  `json:"reqIdentifier"`
	MsgIncr       string `json:"msgIncr"`
	OperationID   string `json:"operationID"`
	ErrCode       int32  `json:"errCode"`
	ErrMsg        string `json:"errMsg"`
	Data          []byte `json:"data"`
}

// argsValidate - 参数校验
func (ws *WServer) argsValidate(m *Req, r int32, operationID string) (isPass bool, errCode int32, errMsg string, returnData interface{}) {
	switch r {
	case constant.WSGetNewestSeq:
		data := open_im_sdk.GetMaxAndMinSeqReq{}
		if err := proto.Unmarshal(m.Data, &data); err != nil {
			log.Error(operationID, "Decode Data struct  err", err.Error(), r)
			return false, 203, err.Error(), nil
		}
		if err := validate.Struct(data); err != nil {
			log.Error(operationID, "data args validate  err", err.Error(), r)
			return false, 204, err.Error(), nil

		}
		return true, 0, "", data
	case constant.WSSendMsg:
		data := open_im_sdk.MsgData{}
		if err := proto.Unmarshal(m.Data, &data); err != nil {
			log.Error(operationID, "Decode Data struct  err", err.Error(), r)
			return false, 203, err.Error(), nil
		}
		if err := validate.Struct(data); err != nil {
			log.Error(operationID, "data args validate  err", err.Error(), r)
			return false, 204, err.Error(), nil
		}
		return true, 0, "", data
	default:
	}
	return false, 204, "args err", nil
}

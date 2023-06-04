package call_back_struct

// 添加好友前回调
type CallbackBeforeAddFriendReq struct {
	CallbackCommand string `json:"callbackCommand"`
	FromUserID      string `json:"fromUserID" `
	ToUserID        string `json:"toUserID"`
	ReqMsg          string `json:"reqMsg"` // 验证信息
	OperationID     string `json:"operationID"`
}

type CallbackBeforeAddFriendResp struct {
	*CommonCallbackResp
}

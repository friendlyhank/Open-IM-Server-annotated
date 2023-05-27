package call_back_struct

// CommonCallbackResp - 回调结果信息
type CommonCallbackResp struct {
	ActionCode  int    `json:"actionCode"`
	ErrCode     int    `json:"errCode"`
	ErrMsg      string `json:"errMsg"`
	OperationID string `json:"operationID"`
}

// UserStatusBaseCallback - 回调状态
type UserStatusBaseCallback struct {
	CallbackCommand string `json:"callbackCommand"` // 回调指令
	OperationID     string `json:"operationID"`
	PlatformID      int32  `json:"platformID"`
	Platform        string `json:"platform"`
}

// 回调状态请求
type UserStatusCallbackReq struct {
	UserStatusBaseCallback
	UserID string `json:"userID"`
}

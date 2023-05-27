package call_back_struct

// CallbackUserOnlineReq - 用户在线回调请求
type CallbackUserOnlineReq struct {
	UserStatusCallbackReq
	Token  string `json:"token"`
	Seq    int    `json:"seq"`
	ConnID string `json:"connID"`
}

// CallbackUserOnlineResp - 用户在线回调结果返回
type CallbackUserOnlineResp struct {
	*CommonCallbackResp
}

// CallbackUserOfflineReq - 用户离线回调请求
type CallbackUserOfflineReq struct {
	UserStatusCallbackReq
	Seq    int    `json:"seq"`
	ConnID string `json:"connID"`
}

// CallbackUserOfflineResp - 用户离线回调结果
type CallbackUserOfflineResp struct {
	*CommonCallbackResp
}

// CallbackUserKickOffReq - 用户踢出回调请求
type CallbackUserKickOffReq struct {
	UserStatusCallbackReq
	Seq int `json:"seq"`
}

// CallbackUserKickOffResp - 用户回调踢出返回
type CallbackUserKickOffResp struct {
	*CommonCallbackResp
}

package base_info

type ParamsCommFriend struct {
	OperationID string `json:"operationID" binding:"required"`
	ToUserID    string `json:"toUserID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}

// 添加好友请求
type AddFriendReq struct {
	ParamsCommFriend
	ReqMsg string `json:"reqMsg"`
}

// 添加好友响应
type AddFriendResp struct {
	CommResp
}

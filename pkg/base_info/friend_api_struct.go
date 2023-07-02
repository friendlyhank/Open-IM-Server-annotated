package base_info

import open_im_sdk "Open_IM/pkg/proto/sdk_ws"

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

// GetFriendApplyListReq - 获取好友申请列表
type GetFriendApplyListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"` // 查询好友列表的用户id
}

type GetFriendApplyListResp struct {
	CommResp
	FriendRequestList []*open_im_sdk.FriendRequest `json:"-"` // 好友申请列表
	Data              []map[string]interface{}     `json:"data" swaggerignore:"true"`
}

// GetSelfApplyListReq - 获取自己的好友申请列表
type GetSelfApplyListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}
type GetSelfApplyListResp struct {
	CommResp
	FriendRequestList []*open_im_sdk.FriendRequest `json:"-"`
	Data              []map[string]interface{}     `json:"data" swaggerignore:"true"`
}

package base_info

import open_im_sdk "Open_IM/pkg/proto/sdk_ws"

// GetUsersInfoReq - 获取用户信息
type GetUsersInfoReq struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required"`
}

// GetUsersInfoResp - 获取用户信息响应接口
type GetUsersInfoResp struct {
	CommResp
	UserInfoList []*open_im_sdk.PublicUserInfo `json:"-"`
	Data         []map[string]interface{}      `json:"data" swaggerignore:"true"`
}

// 获取自己用户信息请求接口
type GetSelfUserInfoReq struct {
	OperationID string `json:"operationID" binding:"required"`
	UserID      string `json:"userID" binding:"required"`
}

// 获取自己用户信息响应接口
type GetSelfUserInfoResp struct {
	CommResp
	UserInfo *open_im_sdk.UserInfo  `json:"-"`
	Data     map[string]interface{} `json:"data" swaggerignore:"true"`
}

package base_info

import open_im_sdk "Open_IM/pkg/proto/sdk_ws"

// 获取用户信息请求接口
type GetSelfUserInfoReq struct {
	OperationID string `json:"operationID" binding:"required"`
	UserID      string `json:"userID" binding:"required"`
}

// 获取用户信息响应接口
type GetSelfUserInfoResp struct {
	CommResp
	UserInfo *open_im_sdk.UserInfo  `json:"-"`
	Data     map[string]interface{} `json:"data" swaggerignore:"true"`
}

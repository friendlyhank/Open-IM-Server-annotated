package base_info

// 用户注册api请求
type UserRegisterReq struct {
	Platform int32 `json:"platform" binding:"required,min=1,max=12"` // 平台
	ApiUserInfo
	OperationID string `json:"operationID" binding:"required"`
}

// 用户token信息
type UserTokenInfo struct {
	UserID      string `json:"userID"`
	Token       string `json:"token"`
	ExpiredTime int64  `json:"expiredTime"`
}

// UserRegisterResp - 用户注册api响应
type UserRegisterResp struct {
	CommResp
	UserToken UserTokenInfo `json:"data"`
}

// UserTokenReq - 用户登录请求
type UserTokenReq struct {
	Secret      string `json:"secret" binding:"required,max=32"`
	Platform    int32  `json:"platform" binding:"required,min=1,max=12"`
	UserID      string `json:"userID" binding:"required,min=1,max=64"`
	LoginIp     string `json:"loginIp"`
	OperationID string `json:"operationID" binding:"required"`
}

type UserTokenResp struct {
	CommResp
	UserToken UserTokenInfo `json:"data"`
}

// ParseTokenReq - 解析用户token请求
type ParseTokenReq struct {
	OperationID string `json:"operationID" binding:"required"`
}

type ExpireTime struct {
	ExpireTimeSeconds uint32 `json:"expireTimeSeconds" `
}

// ParseTokenResp - 解析用户token响应
type ParseTokenResp struct {
	CommResp
	Data       map[string]interface{} `json:"data" swaggerignore:"true"`
	ExpireTime ExpireTime             `json:"-"`
}

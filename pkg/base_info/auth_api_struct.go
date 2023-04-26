package base_info

// 用户注册api请求
type UserRegisterReq struct {
	Platform int32 `json:"platform" binding:"required,min=1,max=12"` // 平台
	ApiUserInfo
	OperationID string `json:"operationID" binding:"required"`
}

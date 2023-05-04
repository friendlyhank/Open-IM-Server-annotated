package base_info

type ApiUserInfo struct {
	UserID      string `json:"userID" binding:"required,min=1,max=64" swaggo:"true,用户ID,"`
	Nickname    string `json:"nickname" binding:"omitempty,min=1,max=64" swaggo:"true,my id,19"`
	FaceURL     string `json:"faceURL" binding:"omitempty,max=1024"`   // 头像
	Gender      int32  `json:"gender" binding:"omitempty,oneof=0 1 2"` // 性别
	PhoneNumber string `json:"phoneNumber" binding:"omitempty,max=32"` // 手机号
	Birth       uint32 `json:"birth" binding:"omitempty"`              // 生日
	Email       string `json:"email" binding:"omitempty,max=64"`
	CreateTime  int64  `json:"createTime"`
}

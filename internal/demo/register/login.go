package register

type ParamsLogin struct {
	UserID      string `json:"userID"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
	Platform    int32  `json:"platform"`
	OperationID string `json:"operationID" binding:"required"`
	AreaCode    string `json:"areaCode"`
}

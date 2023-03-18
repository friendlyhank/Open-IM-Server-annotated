package gate

type Req struct {
	ReqIdentifier int32  `json:"reqIdentifier" validate:"required"` // 请求标识，对应socket协议
	Token         string `json:"token" `
	SendID        string `json:"sendID" validate:"required"` // 发送id
	OperationID   string `json:"operationID" validate:"required"`
	Data          []byte `json:"data"`
}

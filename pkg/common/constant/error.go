package constant

import "errors"

// 错误处理

// key = errCode, string = errMsg
type ErrInfo struct {
	ErrCode int32
	ErrMsg  string
}

var (
	OK        = ErrInfo{0, ""}
	ErrServer = ErrInfo{500, "server error"} // 服务器错误

	ErrArgs = ErrInfo{ErrCode: 803, ErrMsg: ArgsMsg.Error()} // 参数错误
)

var (
	ArgsMsg = errors.New("args failed")
)

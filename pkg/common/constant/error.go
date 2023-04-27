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
	ErrServer = ErrInfo{500, "server error"}                   // 服务器错误
	ErrDB     = ErrInfo{ErrCode: 802, ErrMsg: DBMsg.Error()}   // 数据库mysql错误
	ErrArgs   = ErrInfo{ErrCode: 803, ErrMsg: ArgsMsg.Error()} // 参数错误
)

var (
	DBMsg   = errors.New("db failed")
	ArgsMsg = errors.New("args failed")
)

const (
	RegisterLimit   = 10012 // 用户注册限制
	InvitationError = 10014 // 邀请码错误
)

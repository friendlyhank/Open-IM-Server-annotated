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

	ErrTokenExpired             = ErrInfo{701, TokenExpiredMsg.Error()}             // token过期
	ErrTokenInvalid             = ErrInfo{702, TokenInvalidMsg.Error()}             // token失效
	ErrTokenMalformed           = ErrInfo{703, TokenMalformedMsg.Error()}           // token不正确
	ErrTokenNotValidYet         = ErrInfo{704, TokenNotValidYetMsg.Error()}         // token未激活
	ErrTokenUnknown             = ErrInfo{705, TokenUnknownMsg.Error()}             // 未知token
	ErrTokenKicked              = ErrInfo{706, TokenUserKickedMsg.Error()}          // 被踢出token
	ErrTokenDifferentPlatformID = ErrInfo{707, TokenDifferentPlatformIDMsg.Error()} // token对应平台不同错误
	ErrTokenDifferentUserID     = ErrInfo{708, TokenDifferentUserIDMsg.Error()}     // token对应用户id不同错误

	ErrDB   = ErrInfo{ErrCode: 802, ErrMsg: DBMsg.Error()}   // 数据库mysql错误
	ErrArgs = ErrInfo{ErrCode: 803, ErrMsg: ArgsMsg.Error()} // 参数错误
)

var (
	TokenExpiredMsg             = errors.New("token is timed out, please log in again")
	TokenInvalidMsg             = errors.New("token has been invalidated")
	TokenNotValidYetMsg         = errors.New("token not active yet")
	TokenMalformedMsg           = errors.New("that's not even a token")
	TokenUnknownMsg             = errors.New("couldn't handle this token")
	TokenUserKickedMsg          = errors.New("user has been kicked")
	TokenDifferentPlatformIDMsg = errors.New("different platformID")
	TokenDifferentUserIDMsg     = errors.New("different userID")

	DBMsg   = errors.New("db failed")
	ArgsMsg = errors.New("args failed")
)

const (
	RegisterLimit   = 10012 // 用户注册限制
	InvitationError = 10014 // 邀请码错误
)

func (e ErrInfo) Error() string {
	return e.ErrMsg
}

func (e *ErrInfo) Code() int32 {
	return e.ErrCode
}

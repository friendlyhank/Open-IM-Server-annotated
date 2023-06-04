package token_verify

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	commonDB "Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	go_redis "github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

//var (
//	TokenExpired     = errors.New("token is timed out, please log in again")
//	TokenInvalid     = errors.New("token has been invalidated")
//	TokenNotValidYet = errors.New("token not active yet")
//	TokenMalformed   = errors.New("that's not even a token")
//	TokenUnknown     = errors.New("couldn't handle this token")
//)

/*
 * github.com/golang-jwt/jwt/v4使用第三方token校验工具
 * token校验组件
 */

type Claims struct {
	UID      string
	Platform string //login platform
	jwt.RegisteredClaims
}

func BuildClaims(uid, platform string, ttl int64) Claims {
	now := time.Now()
	return Claims{
		UID:      uid,
		Platform: platform,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(ttl*24) * time.Hour)), //Expiration time
		}}
}

// CreateToken - 创建token
func CreateToken(userID string, platformID int) (string, int64, error) {
	claims := BuildClaims(userID, constant.PlatformIDToName(platformID), config.Config.TokenPolicy.AccessExpire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.TokenPolicy.AccessSecret))
	if err != nil {
		return "", 0, err
	}
	//remove Invalid token
	m, err := commonDB.DB.GetTokenMapByUidPid(userID, constant.PlatformIDToName(platformID))
	if err != nil && err != go_redis.Nil {
		return "", 0, err
	}
	// 删除过期的tokens
	var deleteTokenKey []string
	for k, v := range m {
		_, err = GetClaimFromToken(k)
		if err != nil || v != constant.NormalToken {
			deleteTokenKey = append(deleteTokenKey, k)
		}
	}
	if len(deleteTokenKey) != 0 {
		err = commonDB.DB.DeleteTokenByUidPid(userID, platformID, deleteTokenKey)
		if err != nil {
			return "", 0, err
		}
	}
	err = commonDB.DB.AddTokenFlag(userID, platformID, tokenString, constant.NormalToken)
	if err != nil {
		return "", 0, err
	}
	return tokenString, claims.ExpiresAt.Time.Unix(), err
}

// secret - 常用的设置配置信息的写法
func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.TokenPolicy.AccessSecret), nil
	}
}

// GetClaimFromToken - 根据token获取注册声明信息
func GetClaimFromToken(tokensString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokensString, &Claims{}, secret())
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, utils.Wrap(constant.ErrTokenMalformed, "")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, utils.Wrap(constant.ErrTokenExpired, "")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, utils.Wrap(constant.ErrTokenNotValidYet, "")
			} else {
				return nil, utils.Wrap(constant.ErrTokenUnknown, "")
			}
		} else {
			return nil, utils.Wrap(constant.ErrTokenNotValidYet, "")
		}
	} else {
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return claims, nil
		}
		return nil, utils.Wrap(constant.ErrTokenNotValidYet, "")
	}
}

// 从token中获取用户id
func GetUserIDFromToken(token string, operationID string) (bool, string, string) {
	claims, err := ParseToken(token, operationID)
	if err != nil {
		log.Error(operationID, "ParseToken failed, ", err.Error(), token)
		return false, "", err.Error()
	}
	log.Debug(operationID, "token claims.ExpiresAt.Second() ", claims.ExpiresAt.Unix())
	return true, claims.UID, ""
}

// ParseToken - 解析token
func ParseToken(tokensString, operationID string) (claims *Claims, err error) {
	claims, err = GetClaimFromToken(tokensString)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}

	m, err := commonDB.DB.GetTokenMapByUidPid(claims.UID, claims.Platform)
	if err != nil {
		log.NewError(operationID, "get token from redis err", err.Error(), tokensString)
		return nil, utils.Wrap(constant.ErrTokenInvalid, "get token from redis err")
	}
	if m == nil {
		log.NewError(operationID, "get token from redis err, not in redis ", "m is nil", tokensString)
		return nil, utils.Wrap(constant.ErrTokenInvalid, "get token from redis err")
	}
	if v, ok := m[tokensString]; ok {
		switch v {
		case constant.NormalToken:
			log.NewDebug(operationID, "this is normal return", claims)
			return claims, nil
		case constant.KickedToken:
			log.Error(operationID, "this token has been kicked by other same terminal ", constant.ErrTokenKicked)
			return nil, utils.Wrap(constant.ErrTokenKicked, "this token has been kicked by other same terminal ")
		default:
			return nil, utils.Wrap(constant.ErrTokenUnknown, "")
		}
	}
	log.NewError(operationID, "redis token map not find", constant.ErrTokenUnknown)
	return nil, utils.Wrap(constant.ErrTokenUnknown, "redis token map not find")
}

// WsVerifyToken - 校验websocket的token信息
func WsVerifyToken(token, uid string, platformID string, operationID string) (bool, error, string) {
	argMsg := "args: token: " + token + " operationID: " + operationID + " userID: " + uid + " platformID: " + constant.PlatformIDToName(utils.StringToInt(platformID))
	claims, err := ParseToken(token, operationID)
	if err != nil {
		//if errors.Is(err, constant.ErrTokenUnknown) {
		//	errMsg := "ParseToken failed ErrTokenUnknown " + err.Error()
		//	log.Error(operationID, errMsg)
		//}
		//e := errors.Unwrap(err)
		//if errors.Is(e, constant.ErrTokenUnknown) {
		//	errMsg := "ParseToken failed ErrTokenUnknown " + e.Error()
		//	log.Error(operationID, errMsg)
		//}

		errMsg := "parse token err " + err.Error() + argMsg
		return false, utils.Wrap(err, errMsg), errMsg
	}
	if claims.UID != uid {
		errMsg := " uid is not same to token uid " + argMsg + " claims.UID: " + claims.UID
		return false, utils.Wrap(constant.ErrTokenDifferentUserID, errMsg), errMsg
	}
	if claims.Platform != constant.PlatformIDToName(utils.StringToInt(platformID)) {
		errMsg := " platform is not same to token platform " + argMsg + " claims platformID: " + claims.Platform
		return false, utils.Wrap(constant.ErrTokenDifferentPlatformID, errMsg), errMsg
	}
	log.NewDebug(operationID, utils.GetSelfFuncName(), " check ok ", claims.UID, uid, claims.Platform)
	return true, nil, ""
}

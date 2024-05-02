package controller

import (
	"context"

	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/authverify"

	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/db/cache"

	"github.com/OpenIMSDK/protocol/constant"

	"github.com/OpenIMSDK/tools/errs"

	"github.com/golang-jwt/jwt/v4"

	"github.com/OpenIMSDK/tools/tokenverify"

	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/config"
)

type AuthDatabase interface {
	// If the result is empty, no error is returned.
	GetTokensWithoutError(ctx context.Context, userID string, platformID int) (map[string]int, error)
	// Create token
	CreateToken(ctx context.Context, userID string, platformID int) (string, error)
}

type authDatabase struct {
	cache        cache.MsgModel // 这个分层感觉很有问题，token写到msg
	accessSecret string         // token签名生成密钥
	accessExpire int64          //  接受token过期时间(天)
	config       *config.GlobalConfig
}

func NewAuthDatabase(cache cache.MsgModel, accessSecret string, accessExpire int64, config *config.GlobalConfig) AuthDatabase {
	return &authDatabase{cache: cache, accessSecret: accessSecret, accessExpire: accessExpire, config: config}
}

func (a authDatabase) GetTokensWithoutError(ctx context.Context, userID string, platformID int) (map[string]int, error) {
	//TODO implement me
	panic("implement me")
}

// CreateToken - 创建token
func (a authDatabase) CreateToken(ctx context.Context, userID string, platformID int) (string, error) {
	// 获取用户所有tokens
	tokens, err := a.cache.GetTokensWithoutError(ctx, userID, platformID)
	if err != nil {
		return "", err
	}
	var deleteTokenKey []string
	for k, v := range tokens {
		_, err = tokenverify.GetClaimFromToken(k, authverify.Secret(a.config.Secret))
		if err != nil || v != constant.NormalToken {
			deleteTokenKey = append(deleteTokenKey, k)
		}
	}
	if len(deleteTokenKey) != 0 {
		err = a.cache.DeleteTokenByUidPid(ctx, userID, platformID, deleteTokenKey)
		if err != nil {
			return "", err
		}
	}

	// 生成tokn
	claims := tokenverify.BuildClaims(userID, platformID, a.accessExpire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.accessSecret))
	if err != nil {
		return "", errs.Wrap(err, "token.SignedString")
	}
	return tokenString, a.cache.AddTokenFlag(ctx, userID, platformID, tokenString, constant.NormalToken)
}

package cache

import (
	"context"
	"time"

	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/cachekey"

	"github.com/dtm-labs/rockscache"
	"github.com/redis/go-redis/v9"

	relationtb "github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/db/table/relation"
)

/*
 * 用户缓存
 */
const (
	userExpireTime = time.Second * 60 * 60 * 12 //  用户过期时间
)

type UserCache interface {
	// GetUserInfo - 从缓存获取用户信息
	GetUserInfo(ctx context.Context, userID string) (userInfo *relationtb.UserModel, err error)
	// GetUsersInfo - 批量从缓存获取用户信息
	GetUsersInfo(ctx context.Context, userIDs []string) ([]*relationtb.UserModel, error)
}

type UserCacheRedis struct {
	rdb        redis.UniversalClient
	userDB     relationtb.UserModelInterface
	expireTime time.Duration
	rcClient   *rockscache.Client
}

func NewUserCacheRedis(
	rdb redis.UniversalClient,
	userDB relationtb.UserModelInterface,
	options rockscache.Options,
) UserCache {
	rcClient := rockscache.NewClient(rdb, options)
	return &UserCacheRedis{
		rdb:        rdb,
		userDB:     userDB,
		expireTime: userExpireTime,
		rcClient:   rcClient,
	}
}

func (u *UserCacheRedis) getUserInfoKey(userID string) string {
	return cachekey.GetUserInfoKey(userID)
}

// GetUserInfo - 从缓存获取用户信息
func (u UserCacheRedis) GetUserInfo(ctx context.Context, userID string) (userInfo *relationtb.UserModel, err error) {
	return getCache(ctx, u.rcClient, u.getUserInfoKey(userID), u.expireTime, func(ctx context.Context) (*relationtb.UserModel, error) {
		return u.userDB.Take(ctx, userID)
	})
}

func (u UserCacheRedis) GetUsersInfo(ctx context.Context, userIDs []string) ([]*relationtb.UserModel, error) {
	return batchGetCache2(ctx, u.rcClient, u.expireTime, userIDs, func(userID string) string {
		return u.getUserInfoKey(userID)
	}, func(ctx context.Context, userID string) (*relationtb.UserModel, error) {
		return u.userDB.Take(ctx, userID)
	})
}

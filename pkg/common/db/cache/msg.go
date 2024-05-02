package cache

import (
	"context"

	"github.com/OpenIMSDK/tools/utils"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/errs"

	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/config"
	"github.com/redis/go-redis/v9"
)

const (
	uidPidToken = "UID_PID_TOKEN_STATUS:" // token缓存
)

type MsgModel interface {
	// AddTokenFlag - 添加token标记到缓存
	AddTokenFlag(ctx context.Context, userID string, platformID int, token string, flag int) error
	// GetTokensWithoutError - 获取token信息带上错误
	GetTokensWithoutError(ctx context.Context, userID string, platformID int) (map[string]int, error)
	// DeleteTokenByUidPid - 删除token信息
	DeleteTokenByUidPid(ctx context.Context, userID string, platformID int, fields []string) error
}

func NewMsgCacheModel(client redis.UniversalClient, config *config.GlobalConfig) MsgModel {
	return &msgCache{rdb: client, config: config}
}

type msgCache struct {
	rdb    redis.UniversalClient
	config *config.GlobalConfig
}

// GetTokensWithoutError - 获取token信息带上错误
func (c *msgCache) GetTokensWithoutError(ctx context.Context, userID string, platformID int) (map[string]int, error) {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	m, err := c.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	mm := make(map[string]int)
	for k, v := range m {
		mm[k] = utils.StringToInt(v)
	}

	return mm, nil
}

// AddTokenFlag - 添加token信息到缓存
func (c *msgCache) AddTokenFlag(ctx context.Context, userID string, platformID int, token string, flag int) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	return errs.Wrap(c.rdb.HSet(ctx, key, token, flag).Err())
}

// DeleteTokenByUidPid - 删除token信息
func (c *msgCache) DeleteTokenByUidPid(ctx context.Context, userID string, platform int, fields []string) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platform)

	return errs.Wrap(c.rdb.HDel(ctx, key, fields...).Err())
}

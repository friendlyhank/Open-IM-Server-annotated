package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/dtm-labs/rockscache"
	"time"
)

// GetDefaultOpt - 获取缓存的默认配置
func GetDefaultOpt() rockscache.Options {
	opts := rockscache.NewDefaultOptions()
	opts.StrongConsistency = true     // 强一致性型
	opts.RandomExpireAdjustment = 0.2 // 为了防止雪崩，过期的矫正指数为0.2

	return opts
}

// getCache - 获取缓存信息
func getCache[T any](ctx context.Context, rcClient *rockscache.Client, key string, expire time.Duration, fn func(ctx context.Context) (T, error)) (T, error) {
	var t T
	var write bool
	v, err := rcClient.Fetch2(ctx, key, expire, func() (s string, err error) {
		t, err = fn(ctx)
		if err != nil {
			return "", err
		}
		bs, err := json.Marshal(t)
		if err != nil {
			return "", errs.Wrap(err, "marshal failed")
		}
		write = true

		return string(bs), nil
	})
	if err != nil {
		return t, errs.Wrap(err)
	}
	if write {
		return t, nil
	}
	if v == "" {
		return t, errs.ErrRecordNotFound.Wrap("cache is not found")
	}
	err = json.Unmarshal([]byte(v), &t)
	if err != nil {
		errInfo := fmt.Sprintf("cache json.Unmarshal failed, key:%s, value:%s, expire:%s", key, v, expire)
		return t, errs.Wrap(err, errInfo)
	}

	return t, nil
}

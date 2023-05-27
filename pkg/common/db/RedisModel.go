package db

import (
	"Open_IM/pkg/common/constant"
	log2 "Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"context"
)

/*
 * redis设置
 */

const (
	uidPidToken = "UID_PID_TOKEN_STATUS:" // 用户token 设置
)

// Store userid and platform class to redis - 设置token
func (d *DataBases) AddTokenFlag(userID string, platformID int, token string, flag int) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	log2.NewDebug("", "add token key is ", key)
	return d.RDB.HSet(context.Background(), key, token, flag).Err()
}

// GetTokenMapByUidPid - 获取使用map存储token
func (d *DataBases) GetTokenMapByUidPid(userID, platformID string) (map[string]int, error) {
	key := uidPidToken + userID + ":" + platformID
	log2.NewDebug("", "get token key is ", key)
	m, err := d.RDB.HGetAll(context.Background(), key).Result()
	mm := make(map[string]int)
	for k, v := range m {
		mm[k] = utils.StringToInt(v)
	}
	return mm, err
}

// SetTokenMapByUidPid - 设置map存储token信息
func (d *DataBases) SetTokenMapByUidPid(userID string, platformID int, m map[string]int) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	mm := make(map[string]interface{})
	for k, v := range m {
		mm[k] = v
	}
	return d.RDB.HSet(context.Background(), key, mm).Err()
}

// DeleteTokenByUidPid - 删除指定的token
func (d *DataBases) DeleteTokenByUidPid(userID string, platformID int, fields []string) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	return d.RDB.HDel(context.Background(), key, fields...).Err()
}

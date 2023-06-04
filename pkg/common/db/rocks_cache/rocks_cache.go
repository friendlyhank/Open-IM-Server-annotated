package rocks_cache

import (
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/utils"
	"encoding/json"
	"time"
)

/*
 * 锁缓存，可设置强一致性和弱一致性保证缓存的一致
 */

const (
	userInfoCache = "USER_INFO_CACHE:"
)

// GetUserInfoFromCache - 从缓存中获取用户信息
func GetUserInfoFromCache(userID string) (*db.User, error) {
	getUserInfo := func() (string, error) {
		userInfo, err := imdb.GetUserByUserID(userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(userInfo)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	userInfoStr, err := db.DB.Rc.Fetch(userInfoCache+userID, time.Second*30*60, getUserInfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	userInfo := &db.User{}
	err = json.Unmarshal([]byte(userInfoStr), userInfo)
	return userInfo, utils.Wrap(err, "")
}

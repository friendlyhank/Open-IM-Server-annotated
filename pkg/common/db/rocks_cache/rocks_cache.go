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
	userInfoCache       = "USER_INFO_CACHE:"       // 好友缓存
	friendRelationCache = "FRIEND_RELATION_CACHE:" // 好友关系缓存key
)

// 获取某个用户的好友列表
func GetFriendIDListFromCache(userID string) ([]string, error) {
	getFriendIDList := func() (string, error) {
		friendIDList, err := imdb.GetFriendIDListByUserID(userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(friendIDList)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	friendIDListStr, err := db.DB.Rc.Fetch(friendRelationCache+userID, time.Second*30*60, getFriendIDList)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var friendIDList []string
	err = json.Unmarshal([]byte(friendIDListStr), &friendIDList)
	return friendIDList, utils.Wrap(err, "")
}

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

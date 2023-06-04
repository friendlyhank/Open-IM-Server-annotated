package im_mysql_model

import "Open_IM/pkg/common/db"

// GetFriendIDListByUserID - 获取好友id列表
func GetFriendIDListByUserID(OwnerUserID string) ([]string, error) {
	var friendIDList []string
	err := db.DB.MysqlDB.DefaultGormDB().Table("friends").Where("owner_user_id=?", OwnerUserID).Pluck("friend_user_id", &friendIDList).Error
	if err != nil {
		return nil, err
	}
	return friendIDList, nil
}

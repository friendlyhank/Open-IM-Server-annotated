package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"time"
)

// InsertFriendApplication - 添加好友信息请求
func InsertFriendApplication(friendRequest *db.FriendRequest, args map[string]interface{}) error {
	// todo hank 这里没看懂为啥这么写
	if err := db.DB.MysqlDB.DefaultGormDB().Table("friend_requests").Create(friendRequest).Error; err == nil {
		return nil
	}

	//t := dbConn.Debug().Table("friend_requests").Where("from_user_id = ? and to_user_id = ?", friendRequest.FromUserID, friendRequest.ToUserID).Select("*").Updates(*friendRequest)
	//if t.RowsAffected == 0 {
	//	return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	//}
	//return utils.Wrap(t.Error, "")

	friendRequest.CreateTime = time.Now()
	args["create_time"] = friendRequest.CreateTime
	u := db.DB.MysqlDB.DefaultGormDB().Model(friendRequest).Updates(args)
	//u := dbConn.Table("friend_requests").Where("from_user_id=? and to_user_id=?",
	// friendRequest.FromUserID, friendRequest.ToUserID).Update(&friendRequest)
	//u := dbConn.Table("friend_requests").Where("from_user_id=? and to_user_id=?",
	//	friendRequest.FromUserID, friendRequest.ToUserID).Update(&friendRequest)
	if u.RowsAffected != 0 {
		return nil
	}

	if friendRequest.CreateTime.Unix() < 0 {
		friendRequest.CreateTime = time.Now()
	}
	if friendRequest.HandleTime.Unix() < 0 {
		friendRequest.HandleTime = utils.UnixSecondToTime(0)
	}
	err := db.DB.MysqlDB.DefaultGormDB().Table("friend_requests").Create(friendRequest).Error
	if err != nil {
		return err
	}
	return nil
}

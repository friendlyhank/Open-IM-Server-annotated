package im_mysql_model

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"time"
)

// 用户注册
func UserRegister(user db.User) error {
	user.CreateTime = time.Now()
	if user.AppMangerLevel == 0 {
		user.AppMangerLevel = constant.AppOrdinaryUsers
	}
	if user.Birth.Unix() < 0 {
		user.Birth = utils.UnixSecondToTime(0)
	}
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

// GetAllUser - 获取所有用户
func GetAllUser() ([]db.User, error) {
	var userList []db.User
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Find(&userList).Error
	return userList, err
}

// GetUserByUserID - 根据userid获取用户信息
func GetUserByUserID(userID string) (*db.User, error) {
	var user db.User
	// take复合条件的第一条记录
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Where("user_id=?", userID).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUsersByUserIDList(userIDList []string) ([]*db.User, error) {
	var userList []*db.User
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Where("user_id in (?)", userIDList).Find(&userList).Error
	return userList, err
}

func GetUserNameByUserID(userID string) (string, error) {
	var user db.User
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Select("name").Where("user_id=?", userID).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.Nickname, nil
}

func UpdateUserInfo(user db.User) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("users").Where("user_id=?", user.UserID).Updates(&user).Error
}

func UpdateUserInfoByMap(user db.User, m map[string]interface{}) error {
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Where("user_id=?", user.UserID).Updates(m).Error
	return err
}

func SelectAllUserID() ([]string, error) {
	var resultArr []string
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Pluck("user_id", &resultArr).Error
	if err != nil {
		return nil, err
	}
	return resultArr, nil
}

// GetUserByPhoneNumber - 根据userid获取用户信息
func GetUserByPhoneNumber(phoneNumber string) (*db.User, error) {
	var user db.User
	// take复合条件的第一条记录
	err := db.DB.MysqlDB.DefaultGormDB().Table("users").Where("phone_number=?", phoneNumber).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

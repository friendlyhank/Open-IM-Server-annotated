package db

import "time"

/*
 * 对应表结构
 */

// string UserID = 1;
// string Nickname = 2;
// string FaceUrl = 3;
// int32 Gender = 4;
// string PhoneNumber = 5;
// string Birth = 6;
// string Email = 7;
// string Ex = 8;
// string CreateIp = 9;
// int64 CreateTime = 10;
// int32 AppMangerLevel = 11;
// open_im_sdk.User == imdb.User
type User struct {
	UserID           string    `gorm:"column:user_id;primary_key;size:64"` // 用户id
	Nickname         string    `gorm:"column:name;size:255"`               // 昵称
	FaceURL          string    `gorm:"column:face_url;size:255"`           // 头像
	Gender           int32     `gorm:"column:gender"`                      // 性别
	PhoneNumber      string    `gorm:"column:phone_number;size:32"`        // 手机号码
	Birth            time.Time `gorm:"column:birth"`                       // 生日
	Email            string    `gorm:"column:email;size:64"`               // email
	Ex               string    `gorm:"column:ex;size:1024"`
	CreateTime       time.Time `gorm:"column:create_time;index:create_time"` // 创建时间
	AppMangerLevel   int32     `gorm:"column:app_manger_level"`
	GlobalRecvMsgOpt int32     `gorm:"column:global_recv_msg_opt"`

	status int32 `gorm:"column:status"`
}

// 聊天信息日志表
type ChatLog struct {
	ServerMsgID      string    `gorm:"column:server_msg_id;primary_key;type:char(64)" json:"serverMsgID"`
	ClientMsgID      string    `gorm:"column:client_msg_id;type:char(64)" json:"clientMsgID"`
	SendID           string    `gorm:"column:send_id;type:char(64);index:send_id,priority:2" json:"sendID"` // 发送id
	RecvID           string    `gorm:"column:recv_id;type:char(64);index:recv_id,priority:2" json:"recvID"` // 接受id
	SenderPlatformID int32     `gorm:"column:sender_platform_id" json:"senderPlatformID"`                   // 发送平台id
	SenderNickname   string    `gorm:"column:sender_nick_name;type:varchar(255)" json:"senderNickname"`     // 发送昵称
	SenderFaceURL    string    `gorm:"column:sender_face_url;type:varchar(255);" json:"senderFaceURL"`      // 发送投放
	SessionType      int32     `gorm:"column:session_type;index:session_type,priority:2;index:session_type_alone" json:"sessionType"`
	MsgFrom          int32     `gorm:"column:msg_from" json:"msgFrom"`
	ContentType      int32     `gorm:"column:content_type;index:content_type,priority:2;index:content_type_alone" json:"contentType"`
	Content          string    `gorm:"column:content;type:varchar(3000)" json:"content"`
	Status           int32     `gorm:"column:status" json:"status"`
	SendTime         time.Time `gorm:"column:send_time;index:sendTime;index:content_type,priority:1;index:session_type,priority:1;index:recv_id,priority:1;index:send_id,priority:1" json:"sendTime"` // 发送时间
	CreateTime       time.Time `gorm:"column:create_time" json:"createTime"`                                                                                                                          // 创建时间
	Ex               string    `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}

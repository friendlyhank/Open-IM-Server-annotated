package im_mysql_msg_model

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	pbMsg "Open_IM/pkg/proto/msg"
	"Open_IM/pkg/utils"
	"github.com/jinzhu/copier"
)

/*
 * 消息日志
 */

func InsertMessageToChatLog(msg pbMsg.MsgDataToMQ) error {
	chatLog := new(db.ChatLog)
	copier.Copy(chatLog, msg.MsgData)
	switch msg.MsgData.SessionType {
	case constant.SingleChatType:
		chatLog.RecvID = msg.MsgData.RecvID
	}
	chatLog.Content = string(msg.MsgData.Content)
	chatLog.CreateTime = utils.UnixMillSecondToTime(msg.MsgData.CreateTime)
	log.NewDebug("test", "this is ", chatLog)
	return db.DB.MysqlDB.DefaultGormDB().Table("chat_logs").Create(chatLog).Error
}

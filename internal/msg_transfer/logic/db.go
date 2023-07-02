package logic

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	pbMsg "Open_IM/pkg/proto/msg"
	"Open_IM/pkg/utils"
)

// saveUserChatList - 批量处理聊天信息缓存
func saveUserChatList(userID string, msgList []*pbMsg.MsgDataToMQ, operationID string) (error, uint64) {
	log.Info(operationID, utils.GetSelfFuncName(), "args ", userID, len(msgList))
	//return db.DB.BatchInsertChat(userID, msgList, operationID)
	return db.DB.BatchInsertChat2Cache(userID, msgList, operationID)
}

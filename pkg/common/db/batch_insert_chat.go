package db

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	pbMsg "Open_IM/pkg/proto/msg"
	"Open_IM/pkg/utils"
	"errors"
	go_redis "github.com/go-redis/redis/v8"
)

// BatchInsertChat2Cache - 批量插入聊天缓存
func (d *DataBases) BatchInsertChat2Cache(insertID string, msgList []*pbMsg.MsgDataToMQ, operationID string) (error, uint64) {
	newTime := getCurrentTimestampByMill()
	lenList := len(msgList)
	if lenList > GetSingleGocMsgNum() {
		return errors.New("too large"), 0
	}
	if lenList < 1 {
		return errors.New("too short as 0"), 0
	}
	// judge sessionType to get seq 获取当前最大的req
	var currentMaxSeq uint64
	var err error
	if msgList[0].MsgData.SessionType == constant.SuperGroupChatType {
		currentMaxSeq, err = d.GetGroupMaxSeq(insertID)
		log.Debug(operationID, "constant.SuperGroupChatType  lastMaxSeq before add ", currentMaxSeq, "userID ", insertID, err)
	} else {
		currentMaxSeq, err = d.GetUserMaxSeq(insertID)
		log.Debug(operationID, "constant.SingleChatType  lastMaxSeq before add ", currentMaxSeq, "userID ", insertID, err)
	}
	if err != nil && err != go_redis.Nil {
		return utils.Wrap(err, ""), 0
	}

	lastMaxSeq := currentMaxSeq

	for _, m := range msgList {
		// 累加最大的seq
		currentMaxSeq++
		m.MsgData.Seq = uint32(currentMaxSeq)
		log.Debug(operationID, "cache msg node ", m.String(), m.MsgData.ClientMsgID, "userID: ", insertID, "seq: ", currentMaxSeq)
	}
	log.Debug(operationID, "SetMessageToCache ", insertID, len(msgList))
	// todo hank 维护一个消息缓存
	log.Debug(operationID, "batch to redis  cost time ", getCurrentTimestampByMill()-newTime, insertID, len(msgList))
	if msgList[0].MsgData.SessionType == constant.SuperGroupChatType {
		err = d.SetGroupMaxSeq(insertID, currentMaxSeq)
	} else {
		err = d.SetUserMaxSeq(insertID, currentMaxSeq)
	}
	return utils.Wrap(err, ""), lastMaxSeq
}

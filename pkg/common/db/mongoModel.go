package db

import "time"

const singleGocMsgNum = 5000 // 单次缓存消息的最大数

// GetSingleGocMsgNum - 单次缓存消息的最大数
func GetSingleGocMsgNum() int {
	return singleGocMsgNum
}

// getCurrentTimestampByMill - 毫秒时间戳
func getCurrentTimestampByMill() int64 {
	return time.Now().UnixNano() / 1e6
}

package msg_transfer

import (
	kfk "Open_IM/pkg/common/kafka"
)

/*
 * 持久化消息到mongo
 */

type OnlineHistoryMongoConsumerHandler struct {
	msgHandle            map[string]fcb
	historyConsumerGroup *kfk.MConsumerGroup
}

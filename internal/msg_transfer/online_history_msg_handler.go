package msg_transfer

import "github.com/Shopify/sarama"

type fcb func(cMsg *sarama.ConsumerMessage, msgKey string, sess sarama.ConsumerGroupSession)

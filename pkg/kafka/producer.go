package kafka

import (
	"Open_IM/pkg/common/config"
	"github.com/Shopify/sarama"
)

// kafa生产者

type Producer struct {
	topic    string
	addr     []string
	config   *sarama.Config
	producer sarama.SyncProducer
}

// NewKafkaProducer - 初始化kafka生产者
func NewKafkaProducer(addr []string, topic string) *Producer {
	p := Producer{}
	p.config = sarama.NewConfig()             //Instantiate a sarama Config
	p.config.Producer.Return.Successes = true //Whether to enable the successes channel to be notified after the message is sent successfully
	p.config.Producer.Return.Errors = true
	p.config.Producer.RequiredAcks = sarama.WaitForAll        //Set producer Message Reply level 0 1 all 生产者等待所有消息
	p.config.Producer.Partitioner = sarama.NewHashPartitioner //Set the hash-key automatic hash partition. When sending a message, you must specify the key value of the message. If there is no key, the partition will be selected randomly
	if config.Config.Kafka.SASLUserName != "" && config.Config.Kafka.SASLPassword != "" {
		p.config.Net.SASL.Enable = true
		p.config.Net.SASL.User = config.Config.Kafka.SASLUserName
		p.config.Net.SASL.Password = config.Config.Kafka.SASLPassword
	}
	p.addr = addr
	p.topic = topic
	producer, err := sarama.NewSyncProducer(p.addr, p.config) //Initialize the client
	if err != nil {
		panic(err.Error())
		return nil
	}
	p.producer = producer
	return &p
}

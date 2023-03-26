package msg_transfer

var (
	persistentCH PersistentConsumerHandler
)

func Init() {
	persistentCH.Init() // 持久化消息到数据库
}

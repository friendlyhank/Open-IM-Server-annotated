package msg_transfer

var (
	persistentCH PersistentConsumerHandler
)

func Init() {
	persistentCH.Init() // ws2mschat save mysql
}

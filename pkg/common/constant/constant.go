package constant

const (
	//Websocket Protocol  socket协议
	WSSendMsg = 1003 // 发送消息
	WSPushMsg = 2001 // 推送消息

	//SessionType
	SingleChatType = 1 // 单聊消息
	GroupChatType  = 2 // 群聊消息
	//token 用户token相关
	NormalToken = 0 // 普通token

	OnlineStatus  = "online"  // 在线状态
	OfflineStatus = "offline" // 离线状态
)

const (
	AppOrdinaryUsers = 1 // app普通用户
	AppAdmin         = 2 // app管理员
)

const LogFileName = "OpenIM.log" // 初始化日志名称

const CurrentVersion = "v2.3.8-rc0" // 当前版本

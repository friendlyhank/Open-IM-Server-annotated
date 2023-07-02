package constant

const (
	//Websocket Protocol  socket协议
	WSGetNewestSeq  = 1001 // 获取最新的seq序号
	WSSendMsg       = 1003 // 发送消息
	WSPushMsg       = 2001 // 推送消息
	WSKickOnlineMsg = 2002 // 踢出消息
	WsLogoutMsg     = 2003 // 登出消息

	//SessionType
	SingleChatType       = 1 // 单聊消息
	GroupChatType        = 2 // 群聊消息
	SuperGroupChatType   = 3 // 超级群聊消息
	NotificationChatType = 4 // 通知类型聊天信息
	//token 用户token相关
	NormalToken  = 0 // 普通token
	InValidToken = 1 // 无效token
	KickedToken  = 2 // 被踢出
	ExpiredToken = 3 // 已过期

	//MultiTerminalLogin
	DefalutNotKick = 0 // 多端登录默认不互踢
	//Full-end login, but the same end is mutually exclusive 多个端可以登录，但同端会互踢
	AllLoginButSameTermKick = 1 //
	//Only one of the endpoints can log in 只有一个端可以登录
	SingleTerminalLogin = 2
	//The web side can be online at the same time, and the other side can only log in at one end web端不互踢，其他端互踢
	WebAndOther = 3
	//The PC side is mutually exclusive, and the mobile side is mutually exclusive, but the web side can be online at the same time pc和mobile不互踢，其他端互踢
	PcMobileAndWeb = 4
	//The PC terminal can be online at the same time,but other terminal only one of the endpoints can login
	PCAndOther = 5 // pc端不互踢，其他端互踢

	OnlineStatus  = "online"  // 在线状态
	OfflineStatus = "offline" // 离线状态

	//callbackCommand 回调指令
	CallbackUserOnlineCommand      = "callbackUserOnlineCommand"      // 在线回调指令
	CallbackUserOfflineCommand     = "callbackUserOfflineCommand"     // 离线回调指令
	CallbackUserKickOffCommand     = "callbackUserKickOffCommand"     // 踢出回调指令
	CallbackBeforeAddFriendCommand = "callbackBeforeAddFriendCommand" // 添加好友前指令

	//callback actionCode
	ActionAllow     = 0 // 回调失败继续往下操作
	ActionForbidden = 1 // 回调失败不继续往下
)

const (
	AppOrdinaryUsers = 1 // app普通用户
	AppAdmin         = 2 // app管理员
)

const LogFileName = "OpenIM.log" // 初始化日志名称

const StatisticsTimeInterval = 60 // 统计信息时间间隔

const CurrentVersion = "v2.3.8-rc0" // 当前版本

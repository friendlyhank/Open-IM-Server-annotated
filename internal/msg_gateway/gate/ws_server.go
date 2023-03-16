package gate

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type UserConn struct {
	*websocket.Conn             // 连接
	w               *sync.Mutex // 锁
	PlatformID      int32       // 平台id
	IsCompress      bool        // 是否压缩数据
	userID          string      // 用户id
	token           string      // token
	connID          string      // 连接connID
}

type WServer struct {
	wsMaxConnNum int                            // 最大连接数
	wsUpGrader   *websocket.Upgrader            // socket初始化配置
	wsUserToConn map[string]map[int][]*UserConn // 存储用户连接信息
}

// 初始化WServer
func (ws *WServer) onInit(wsPort int) {
	ws.wsAddr = ":" + utils.IntToString(wsPort)
	ws.wsUserToConn = make(map[string]map[int][]*UserConn)
	ws.wsUpGrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }, // 允许跨域
	}
}

func (ws *WServer) run() {
	http.HandleFunc("/", ws.wsHandler)         //Get request from client to handle by wsHandler
	err := http.ListenAndServe(ws.wsAddr, nil) //Start listening
	if err != nil {
		panic("Ws listening err:" + err.Error())
	}
}

func (ws *WServer) wsHandler(w http.ResponseWriter, r *http.Request) {}

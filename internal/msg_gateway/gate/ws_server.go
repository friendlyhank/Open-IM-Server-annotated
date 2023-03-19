package gate

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/utils"
	"bytes"
	"compress/gzip"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// im服务实现

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
	wsAddr       string                         // 地址
	wsMaxConnNum int                            // 最大连接数
	wsUpGrader   *websocket.Upgrader            // socket初始化配置
	wsUserToConn map[string]map[int][]*UserConn // 存储用户连接信息 todo hank 为什么要设置切片
}

// 初始化WServer
func (ws *WServer) onInit(wsPort int) {
	ws.wsAddr = ":" + utils.IntToString(wsPort)
	ws.wsUserToConn = make(map[string]map[int][]*UserConn)
	ws.wsUpGrader = &websocket.Upgrader{
		HandshakeTimeout: time.Duration(config.Config.LongConnSvr.WebsocketTimeOut) * time.Second,
		ReadBufferSize:   config.Config.LongConnSvr.WebsocketMaxMsgLen,
		CheckOrigin:      func(r *http.Request) bool { return true }, // 允许跨域
	}
}

// run - 服务运行
func (ws *WServer) run() {
	http.HandleFunc("/", ws.wsHandler)         //Get request from client to handle by wsHandler
	err := http.ListenAndServe(ws.wsAddr, nil) //Start listening
	if err != nil {
		panic("Ws listening err:" + err.Error())
	}
}

// wsHandler - socket连接处理逻辑
func (ws *WServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// operationID用于链路追踪
	operationID := ""
	if len(query["operationID"]) != 0 {
		operationID = query["operationID"][0]
	} else {
		operationID = utils.OperationIDGenerator()
	}
	if isPass, compression := ws.headerCheck(w, r, operationID); isPass {
		conn, err := ws.wsUpGrader.Upgrade(w, r, nil) //Conn is obtained through the upgraded escalator
		if err != nil {
			return
		} else {
			newConn := &UserConn{conn, new(sync.Mutex), utils.StringToInt32(query["platformID"][0]), compression, query["sendID"][0], query["token"][0], utils.Md5(conn.RemoteAddr().String() + "_" + strconv.Itoa(int(utils.GetCurrentTimestampByMill())))}
			userCount++
			ws.addUserConn(query["sendID"][0], utils.StringToInt(query["platformID"][0]), newConn, query["token"][0], newConn.connID, operationID)
			go ws.readMsg(newConn)
		}
	} else {
	}
}

// readMsg - 读取消息
func (ws *WServer) readMsg(conn *UserConn) {
	for {
		messageType, msg, err := conn.ReadMessage()
		if messageType == websocket.PingMessage {
		}
		if err != nil {
			userCount--
			ws.delUserConn(conn)
			return
		}
		if messageType == websocket.CloseMessage {
			userCount--
			ws.delUserConn(conn)
			return
		}
		// 开启数据压缩
		if conn.IsCompress {
			buff := bytes.NewBuffer(msg)
			reader, err := gzip.NewReader(buff)
			if err != nil {
				continue
			}
			msg, err = ioutil.ReadAll(reader)
			if err != nil {
				continue
			}
			err = reader.Close()
			if err != nil {
			}
		}
		ws.msgParse(conn, msg)
	}
}

// writeMsg - 发送消息
func (ws *WServer) writeMsg(conn *UserConn, a int, msg []byte) error {
	conn.w.Lock()
	defer conn.w.Unlock()
	if conn.IsCompress {
		var buffer bytes.Buffer
		gz := gzip.NewWriter(&buffer)
		if _, err := gz.Write(msg); err != nil {
		}
		if err := gz.Close(); err != nil {
		}
		msg = buffer.Bytes()
	}
	conn.SetWriteDeadline(time.Now().Add(time.Duration(60) * time.Second))
	return conn.WriteMessage(a, msg)
}

// addUserConn - 添加用户连接
func (ws *WServer) addUserConn(uid string, platformID int, conn *UserConn, token string, connID, operationID string) {
	rwLock.Lock()
	defer rwLock.Unlock()

	// 用户上线回调

	if oldConnMap, ok := ws.wsUserToConn[uid]; ok {
		if conns, ok := oldConnMap[platformID]; ok {
			conns = append(conns, conn)
			oldConnMap[platformID] = conns
		} else {
			var conns []*UserConn
			conns = append(conns, conn)
			oldConnMap[platformID] = conns
		}
		ws.wsUserToConn[uid] = oldConnMap
	} else {
		i := make(map[int][]*UserConn)
		var conns []*UserConn
		conns = append(conns, conn)
		i[platformID] = conns
		ws.wsUserToConn[uid] = i
	}
}

// delUserConn - 删除用户连接
func (ws *WServer) delUserConn(conn *UserConn) {
	rwLock.Lock()
	defer rwLock.Unlock()
	platform := int(conn.PlatformID)

	if oldConnMap, ok := ws.wsUserToConn[conn.userID]; ok { // only recycle self conn
		if oldconns, okMap := oldConnMap[platform]; okMap {

			var a []*UserConn

			for _, client := range oldconns {
				if client != conn {
					a = append(a, client)

				}
			}
			if len(a) != 0 {
				oldConnMap[platform] = a
			} else {
				delete(oldConnMap, platform)
			}
		}
		ws.wsUserToConn[conn.userID] = oldConnMap
		if len(oldConnMap) == 0 {
			delete(ws.wsUserToConn, conn.userID)
		}
	}

	err := conn.Close()
	if err != nil {
	}
	if conn.PlatformID == 0 || conn.connID == "" {
	}
	// todo hank 用户下线回调
}

// headerCheck - 头部信息校验
func (ws *WServer) headerCheck(w http.ResponseWriter, r *http.Request, operationID string) (isPass, compression bool) {
	return true, false
}

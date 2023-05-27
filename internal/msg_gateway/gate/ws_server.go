package gate

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/utils"
	"bytes"
	"compress/gzip"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
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
	wsUserToConn map[string]map[int][]*UserConn // 存储用户连接信息
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
	log.Debug(operationID, utils.GetSelfFuncName(), " args: ", query)
	if isPass, compression := ws.headerCheck(w, r, operationID); isPass {
		conn, err := ws.wsUpGrader.Upgrade(w, r, nil) //Conn is obtained through the upgraded escalator
		if err != nil {
			log.Error(operationID, "upgrade http conn err", err.Error(), query)
			return
		} else {
			newConn := &UserConn{conn, new(sync.Mutex), utils.StringToInt32(query["platformID"][0]), compression, query["sendID"][0], query["token"][0], utils.Md5(conn.RemoteAddr().String() + "_" + strconv.Itoa(int(utils.GetCurrentTimestampByMill())))}
			userCount++
			ws.addUserConn(query["sendID"][0], utils.StringToInt(query["platformID"][0]), newConn, query["token"][0], newConn.connID, operationID)
			go ws.readMsg(newConn)
		}
	} else {
		log.Error(operationID, "headerCheck failed ")
	}
}

// readMsg - 读取消息
func (ws *WServer) readMsg(conn *UserConn) {
	for {
		messageType, msg, err := conn.ReadMessage()
		if messageType == websocket.PingMessage {
			log.NewInfo("", "this is a  pingMessage")
		}
		if err != nil {
			log.NewWarn("", "WS ReadMsg error ", " userIP", conn.RemoteAddr().String(), "userUid", "platform", "error", err.Error())
			userCount--
			ws.delUserConn(conn)
			return
		}
		if messageType == websocket.CloseMessage {
			log.NewWarn("", "WS receive error ", " userIP", conn.RemoteAddr().String(), "userUid", "platform", "error", string(msg))
			userCount--
			ws.delUserConn(conn)
			return
		}
		log.NewDebug("", "size", utils.ByteSize(uint64(len(msg))))
		// 开启数据压缩
		if conn.IsCompress {
			buff := bytes.NewBuffer(msg)
			reader, err := gzip.NewReader(buff)
			if err != nil {
				log.NewWarn("", "un gzip read failed")
				continue
			}
			msg, err = ioutil.ReadAll(reader)
			if err != nil {
				log.NewWarn("", "ReadAll failed")
				continue
			}
			err = reader.Close()
			if err != nil {
				log.NewWarn("", "reader close failed")
			}
		}
		ws.msgParse(conn, msg)
	}
}

// SetWriteTimeout - 设置超时写方法
func (ws *WServer) SetWriteTimeout(conn *UserConn, timeout int) {
	conn.w.Lock()
	defer conn.w.Unlock()
	conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
}

// writeMsg - 发送消息
func (ws *WServer) writeMsg(conn *UserConn, a int, msg []byte) error {
	conn.w.Lock()
	defer conn.w.Unlock()
	if conn.IsCompress {
		var buffer bytes.Buffer
		gz := gzip.NewWriter(&buffer)
		if _, err := gz.Write(msg); err != nil {
			return utils.Wrap(err, "")
		}
		if err := gz.Close(); err != nil {
			return utils.Wrap(err, "")
		}
		msg = buffer.Bytes()
	}
	conn.SetWriteDeadline(time.Now().Add(time.Duration(60) * time.Second))
	return conn.WriteMessage(a, msg)
}

// SetWriteTimeoutWriteMsg - 设置写消息超时
func (ws *WServer) SetWriteTimeoutWriteMsg(conn *UserConn, a int, msg []byte, timeout int) error {
	conn.w.Lock()
	defer conn.w.Unlock()
	conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	return conn.WriteMessage(a, msg)
}

// addUserConn - 添加用户连接
func (ws *WServer) addUserConn(uid string, platformID int, conn *UserConn, token string, connID, operationID string) {
	rwLock.Lock()
	defer rwLock.Unlock()
	log.Info(operationID, utils.GetSelfFuncName(), " args: ", uid, platformID, conn, token, "ip: ", conn.RemoteAddr().String())

	// todo hank 用户上线回调
	//todo hank 多端用户踢出逻辑

	if oldConnMap, ok := ws.wsUserToConn[uid]; ok {
		oldConnMap[platformID] = conn
		ws.wsUserToConn[uid] = oldConnMap
		log.Debug(operationID, "user not first come in, add conn ", uid, platformID, conn, oldConnMap)
	} else {
		i := make(map[int]*UserConn)
		i[platformID] = conn
		ws.wsUserToConn[uid] = i
		log.Debug(operationID, "user first come in, new user, conn", uid, platformID, conn, ws.wsUserToConn[uid])
	}
	if oldStringMap, ok := ws.wsConnToUser[conn]; ok {
		oldStringMap[platformID] = uid
		ws.wsConnToUser[conn] = oldStringMap
	} else {
		i := make(map[int]string)
		i[platformID] = uid
		ws.wsConnToUser[conn] = i
	}
	// 计算用户的连接数
	count := 0
	for _, v := range ws.wsUserToConn {
		count = count + len(v)
	}
	log.Debug(operationID, "WS Add operation", "", "wsUser added", ws.wsUserToConn, "connection_uid", uid, "connection_platform", constant.PlatformIDToName(platformID), "online_user_num", len(ws.wsUserToConn), "online_conn_num", count)
}

// delUserConn - 删除用户连接
func (ws *WServer) delUserConn(conn *UserConn) {
	rwLock.Lock()
	defer rwLock.Unlock()
	operationID := utils.OperationIDGenerator()
	var uid string
	var platform int
	if oldStringMap, okg := ws.wsConnToUser[conn]; okg {
		for k, v := range oldStringMap {
			platform = k
			uid = v
		}
		if oldConnMap, ok := ws.wsUserToConn[uid]; ok {
			delete(oldConnMap, platform)
			ws.wsUserToConn[uid] = oldConnMap
			if len(oldConnMap) == 0 {
				delete(ws.wsUserToConn, uid)
			}
			count := 0
			for _, v := range ws.wsUserToConn {
				count = count + len(v)
			}
			log.Debug(operationID, "WS delete operation", "", "wsUser deleted", ws.wsUserToConn, "disconnection_uid", uid, "disconnection_platform", platform, "online_user_num", len(ws.wsUserToConn), "online_conn_num", count)
		} else {
			log.Debug(operationID, "WS delete operation", "", "wsUser deleted", ws.wsUserToConn, "disconnection_uid", uid, "disconnection_platform", platform, "online_user_num", len(ws.wsUserToConn))
		}
		delete(ws.wsConnToUser, conn)
	}
	err := conn.Close()
	if err != nil {
		log.Error(operationID, " close err", "", "uid", uid, "platform", platform)
	}
	if conn.PlatformID == 0 || conn.connID == "" {
		log.NewWarn(operationID, utils.GetSelfFuncName(), "PlatformID or connID is null", conn.PlatformID, conn.connID)
	}
	// todo hank 用户下线回调
}

// getUserAllCons -获取用户所有的连接
func (ws *WServer) getUserAllCons(uid string) map[int]*UserConn {
	rwLock.RLock()
	defer rwLock.RUnlock()
	if connMap, ok := ws.wsUserToConn[uid]; ok {
		newConnMap := make(map[int]*UserConn)
		for k, v := range connMap {
			newConnMap[k] = v
		}
		return newConnMap
	}
	return nil
}

// headerCheck - 头部信息校验
func (ws *WServer) headerCheck(w http.ResponseWriter, r *http.Request, operationID string) (isPass, compression bool) {
	status := http.StatusUnauthorized
	query := r.URL.Query()
	if len(query["token"]) != 0 && len(query["sendID"]) != 0 && len(query["platformID"]) != 0 {
		if ok, err, msg := token_verify.WsVerifyToken(query["token"][0], query["sendID"][0], query["platformID"][0], operationID); !ok {
			if errors.Is(err, constant.ErrTokenExpired) {
				status = int(constant.ErrTokenExpired.ErrCode)
			}
			if errors.Is(err, constant.ErrTokenInvalid) {
				status = int(constant.ErrTokenInvalid.ErrCode)
			}
			if errors.Is(err, constant.ErrTokenMalformed) {
				status = int(constant.ErrTokenMalformed.ErrCode)
			}
			if errors.Is(err, constant.ErrTokenNotValidYet) {
				status = int(constant.ErrTokenNotValidYet.ErrCode)
			}
			if errors.Is(err, constant.ErrTokenUnknown) {
				status = int(constant.ErrTokenUnknown.ErrCode)
			}
			if errors.Is(err, constant.ErrTokenKicked) {
				status = int(constant.ErrTokenKicked.ErrCode)
			}
			if errors.Is(err, constant.ErrTokenDifferentPlatformID) {
				status = int(constant.ErrTokenDifferentPlatformID.ErrCode)
			}
			if errors.Is(err, constant.ErrTokenDifferentUserID) {
				status = int(constant.ErrTokenDifferentUserID.ErrCode)
			}
			//switch errors.Cause(err) {
			//case constant.ErrTokenExpired:
			//	status = int(constant.ErrTokenExpired.ErrCode)
			//case constant.ErrTokenInvalid:
			//	status = int(constant.ErrTokenInvalid.ErrCode)
			//case constant.ErrTokenMalformed:
			//	status = int(constant.ErrTokenMalformed.ErrCode)
			//case constant.ErrTokenNotValidYet:
			//	status = int(constant.ErrTokenNotValidYet.ErrCode)
			//case constant.ErrTokenUnknown:
			//	status = int(constant.ErrTokenUnknown.ErrCode)
			//case constant.ErrTokenKicked:
			//	status = int(constant.ErrTokenKicked.ErrCode)
			//case constant.ErrTokenDifferentPlatformID:
			//	status = int(constant.ErrTokenDifferentPlatformID.ErrCode)
			//case constant.ErrTokenDifferentUserID:
			//	status = int(constant.ErrTokenDifferentUserID.ErrCode)
			//}

			log.Error(operationID, "Token verify failed ", "query ", query, msg, err.Error(), "status: ", status)
			w.Header().Set("Sec-Websocket-Version", "13")
			w.Header().Set("ws_err_msg", err.Error())
			http.Error(w, err.Error(), status)
			return false, false
		} else {
			if r.Header.Get("compression") == "gzip" {
				compression = true
			}
			if len(query["compression"]) != 0 && query["compression"][0] == "gzip" {
				compression = true
			}
			log.Info(operationID, "Connection Authentication Success", "", "token ", query["token"][0], "userID ", query["sendID"][0], "platformID ", query["platformID"][0], "compression", compression)
			return true, compression
		}
	} else {
		status = int(constant.ErrArgs.ErrCode)
		log.Error(operationID, "Args err ", "query ", query)
		w.Header().Set("Sec-Websocket-Version", "13")
		errMsg := "args err, need token, sendID, platformID"
		w.Header().Set("ws_err_msg", errMsg)
		http.Error(w, errMsg, status)
		return false, false
	}
}

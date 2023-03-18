package gate

import (
	"Open_IM/pkg/common/constant"
	"bytes"
	"encoding/gob"
)

// msgParse - 解析消息
func (ws *WServer) msgParse(conn *UserConn, binaryMsg []byte) {
	b := bytes.NewBuffer(binaryMsg)
	m := Req{}
	dec := gob.NewDecoder(b)
	err := dec.Decode(&m)
	if err != nil {
		err = conn.Close()
		if err != nil {
		}
		return
	}
	if m.SendID != conn.userID {
		if err = conn.Close(); err != nil {
			return
		}
	}
	switch m.ReqIdentifier {
	case constant.WSSendMsg:
	}
}

// sendMsgReq - 发送消息请求
func (ws *WServer) sendMsgReq(conn *UserConn, m *Req) {
	sendMsgAllCountLock.Lock()
	sendMsgAllCount++
	sendMsgAllCountLock.Unlock()
}

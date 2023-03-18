package gate

import (
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

}

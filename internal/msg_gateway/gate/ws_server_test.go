package gate

import (
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"testing"
)

func Test_WServerSendMsg(t *testing.T) {
	url := "ws://0.0.0.0:10001?token=123&sendID=383&platformID=1"
	c, res, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	log.Printf("响应:%s", fmt.Sprint(res))
	defer c.Close()
	done := make(chan struct{})

	// 发送消息结构体
	msgData := &open_im_sdk.MsgData{
		SendID: "1",
	}

	d, _ := proto.Marshal(msgData)

	// im请求结构体
	m := &Req{
		ReqIdentifier: 1003,
		Token:         "123",
		SendID:        "383",
		OperationID:   "1",
		MsgIncr:       "1",
		Data:          d,
	}

	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	enc.Encode(m)

	err = c.WriteMessage(websocket.BinaryMessage, b.Bytes())
	if err != nil {
		fmt.Println(err.Error())
	}
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Fatal(err)
			break
		}
		log.Printf("收到消息: %s", message)
	}
	<-done
}

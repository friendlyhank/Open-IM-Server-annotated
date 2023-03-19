package gate

import (
	"fmt"
	"github.com/gorilla/websocket"
	"testing"
	"time"
)

func Test_WServerSocket(t *testing.T) {
	conn, res, err := websocket.DefaultDialer.Dial("ws://0.0.0.0:10001?token=123&sendID=383&platformID=1", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.WriteMessage(websocket.BinaryMessage, []byte("132"))
	fmt.Println(res.Header)
	time.Sleep(20 * time.Second)
}

package gate

import (
	"fmt"
	"github.com/gorilla/websocket"
	"testing"
	"time"
)

func Test_WServerSocket(t *testing.T) {
	conn, _, err := websocket.DefaultDialer.Dial("ws://0.0.0.0:10001?token=123&sendID=383&platformID=1", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn.WriteMessage(websocket.BinaryMessage, []byte("hello"))

	time.Sleep(5 * time.Minute)
}

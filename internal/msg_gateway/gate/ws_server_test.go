package gate

import (
	"fmt"
	"github.com/gorilla/websocket"
	"testing"
)

func Test_WServerSocket(t *testing.T) {
	conn, _, err := websocket.DefaultDialer.Dial("ws://0.0.0.0:10001/", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(conn)
}

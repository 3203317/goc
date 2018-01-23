package msg

import (
	"fmt"
	"github.com/gorilla/websocket"
)

func Login(conn *websocket.Conn, token string) {
	conn.WriteMessage(websocket.BinaryMessage, []byte(token))
	fmt.Println("send token: ", token)
}

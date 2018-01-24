package msg

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

var (
	HEARTBEAT = flag.Duration("HEARTBEAT", 1, "心跳(秒)")
)

func Login(conn *websocket.Conn, token string) {
	if err := conn.WriteMessage(websocket.BinaryMessage, []byte(token)); nil != err {
		log.Fatal(err)
	}

	fmt.Println("send token: ", token)
}

func Heartbeat(conn *websocket.Conn) {
	for {
		time.Sleep(time.Nanosecond * *HEARTBEAT)

		if err := conn.WriteMessage(websocket.BinaryMessage, []byte("['',7,'']")); nil != err {
			log.Fatal(err)
		}
	}
}

func OnMessage(conn *websocket.Conn) {
	for {
		_, data, err := conn.ReadMessage()

		if nil != err {
			conn.Close()
			log.Fatal(err)
			break
		}

		fmt.Println(string(data))
	}
}

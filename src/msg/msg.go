package msg

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

var (
	HEARTBEAT = flag.Duration("HEARTBEAT", 1, "心跳间隔(秒)")
)

func Login(conn *websocket.Conn, token string) {
	if err := conn.WriteMessage(websocket.BinaryMessage, []byte(token)); nil != err {
		log.Fatal(err)
	}

	fmt.Println("send token: ", token)
}

func Heartbeat(conn *websocket.Conn) {
	ticker := time.NewTicker(time.Nanosecond * *HEARTBEAT)

	b := []byte("['',7,'']")

	for {
		select {
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.BinaryMessage, b); nil != err {
				log.Fatal(err)
			}
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

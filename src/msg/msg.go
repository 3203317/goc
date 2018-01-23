package msg

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

func Login(conn *websocket.Conn, token string) {
	if err := conn.WriteMessage(websocket.BinaryMessage, []byte(token)); nil != err {
		log.Fatal(err)
	}

	fmt.Println("send token: ", token)
}

func Heartbeat(conn *websocket.Conn) {
	for {
		time.Sleep(time.Second * 1)

		if err := conn.WriteMessage(websocket.BinaryMessage, []byte("['',7,'']")); nil != err {
			log.Fatal(err)
		}
	}
}

func Test(conn *websocket.Conn) {
	for {
		_, p, err := conn.ReadMessage()

		if nil != err {
			log.Println(err)
			return
		}

		fmt.Printf("rec: %s\n", p)
	}
}

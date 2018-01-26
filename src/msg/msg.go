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

	fmt.Println("send token:", token)
}

func Heartbeat(conn *websocket.Conn, ch_err_code chan int) {
	ticker := time.NewTicker(time.Nanosecond * *HEARTBEAT)

	defer func() {
		ticker.Stop()
	}()

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

func OnMessage(conn *websocket.Conn, ch_read_msg chan []byte, ch_err_code chan int) {
	for {
		_, msg, err := conn.ReadMessage()

		if nil != err {
			werr := err.(*websocket.CloseError)
			ch_err_code <- werr.Code
			break
		}

		ch_read_msg <- msg
	}
}

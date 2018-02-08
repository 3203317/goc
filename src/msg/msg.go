package msg

import (
	"config"
	"encoding/json"
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

func Heartbeat(conn *websocket.Conn, ch_write_msg chan []byte, ch_err chan error) {
	ticker := time.NewTicker(time.Nanosecond * *HEARTBEAT)

	defer func() {
		ticker.Stop()
	}()

	b := []byte("['',7,'']")

	for {
		select {
		case <-ticker.C:

			if err := conn.WriteMessage(websocket.BinaryMessage, b); nil != err {
				ch_err <- err
			}

		case msg := <-ch_write_msg:

			if err := conn.WriteMessage(websocket.BinaryMessage, msg); nil != err {
				ch_err <- err
			}

		}
	}
}

func OnMessage(conn *websocket.Conn, ch_read_msg chan []byte, ch_status chan config.Status) {
	for {
		_, msg, err := conn.ReadMessage()

		if nil != err {
			// werr := err.(*websocket.CloseError)
			// ch_err <- werr.Code
			ch_status <- config.Status{Code: -2, Err: err}
			break
		}

		if nil != msg {
			ch_read_msg <- msg
		}
	}
}

func Process(ch_read_msg, ch_write_msg chan []byte, ch_err chan error) {
	b := []byte("['',2,'']")

	for {
		select {
		case msg := <-ch_read_msg:

			var sb []interface{}

			if err := json.Unmarshal(msg, &sb); nil != err {
				ch_err <- err
				continue
			}

			switch sb[0].(float64) {
			case 1:
				fmt.Println("登陆:", sb)
			case 2:
				fmt.Println("data:", sb)
			case 7:
				fmt.Println("data:", sb)
				ch_write_msg <- b
			}

		case err := <-ch_err:

			// switch code {
			// case websocket.CloseNoStatusReceived:
			// 	fmt.Println("CloseNoStatusReceived:", code)
			// case websocket.CloseAbnormalClosure:
			// 	fmt.Println("CloseAbnormalClosure:", code)
			// case websocket.CloseMessageTooBig:
			// 	fmt.Println("CloseMessageTooBig:", code)
			// default:
			// 	fmt.Println("code:", code)
			// }

			log.Fatal(err)

		}
	}
}

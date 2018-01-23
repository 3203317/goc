package msg

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
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

func Test(c *websocket.Conn) {
	for {

		_, p, err := c.NextReader()

		if err != nil {
			c.Close()
			break
		}

		aaa, _ := readFrom(p, 128)

		fmt.Println(string(aaa))

		// p1, _ := ioutil.ReadAll(p)

		// fmt.Println(string(p1))
	}
}

func readFrom(reader io.Reader, num int) ([]byte, error) {
	p := make([]byte, num)
	n, err := reader.Read(p)
	if n > 0 {
		return p[:n], nil
	}
	return p, err
}

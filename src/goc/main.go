package main

import (
	"config"
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	UUID "github.com/snluu/uuid"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

var (
	server_addr    = flag.String("server_addr", "47.104.99.102:9988", "前置机地址")
	redis_addr     = flag.String("redis_addr", "47.104.99.102:6379", "Redis地址")
	redis_pwd      = flag.String("redis_pwd", "shuoleniyebudong", "Redis密码")
	redis_sha_auth = flag.String("redis_sha_auth", "a0ad12f31d7de75a5153bdff954caf5bc15b9501", "Redis授权码")
)

func def(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func main() {
	flag.Parse()

	config.LoadCfg()
	fmt.Println("Hello, GO!")

	token := connRedis()

	fmt.Println(token, reflect.TypeOf(token))

	// http.HandleFunc("/", def)
	// err := http.ListenAndServe(":80", nil)
	// if nil != err {
	// 	log.Fatal(err)
	// }

	// runTcpCli(token)

	u := url.URL{Scheme: "ws", Host: *server_addr, Path: "/"}
	var dialer *websocket.Dialer

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	go timeWriter2(conn, token)
	go timeWriter(conn, token)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}

		fmt.Printf("received: %s\n", message)
	}
}

func timeWriter2(conn *websocket.Conn, token string) {
	conn.WriteMessage(websocket.BinaryMessage, []byte(token))
}

func timeWriter(conn *websocket.Conn, token string) {
	for {
		time.Sleep(time.Second * 2)
		// conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
		conn.WriteMessage(websocket.BinaryMessage, []byte("['',7,'']"))
		conn.WriteMessage(websocket.BinaryMessage, []byte("['',2,'']"))
	}
}

func runTcpCli(token string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", *server_addr)

	if nil != err {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	lens, err := conn.Write([]byte(token))
	if nil != err {
		log.Fatal(err)
	}

	fmt.Println(lens)

	for {
	}
}

func connRedis() string {
	client := redis.NewClient(&redis.Options{
		Addr:     *redis_addr,
		Password: *redis_pwd,
		DB:       1,
	})

	defer client.Close()

	// pong, err := client.Ping().Result()
	// if nil != err {
	// 	log.Fatal(err)
	// }

	// log.Println(pong)

	uuid := strings.Replace(UUID.Rand().Hex(), "-", "", -1)

	_token, err := client.EvalSha(*redis_sha_auth, []string{"1", "1", "backend_1", uuid}, 5, 68, "BACK").Result()
	if nil != err {
		log.Fatal(err)
	}

	token, _ := _token.(string)
	return token
}

package main

import (
	"config"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	UUID "github.com/snluu/uuid"
	"io"
	"log"
	"msg"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

var (
	SERVER_ADDR    = flag.String("SERVER_ADDR", "47.104.99.102:9988", "前置机地址")
	REDIS_ADDR     = flag.String("REDIS_ADDR", "47.104.99.102:6379", "Redis地址")
	REDIS_PWD      = flag.String("REDIS_PWD", "shuoleniyebudong", "Redis密码")
	REDIS_SHA_AUTH = flag.String("REDIS_SHA_AUTH", "a0ad12f31d7de75a5153bdff954caf5bc15b9501", "Redis授权码")
)

var (
	ch_read_msg  = make(chan []byte)
	ch_write_msg = make(chan []byte)
	ch_err_code  = make(chan int)
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

	runWsCli(token)
}

func runWsCli(token string) {
	u := url.URL{Scheme: "ws", Host: *SERVER_ADDR, Path: "/"}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if nil != err {
		log.Fatal(err)
	}

	defer conn.Close()

	conn.EnableWriteCompression(true)

	go msg.OnMessage(conn, ch_read_msg, ch_err_code)

	msg.Login(conn, token)

	go msg.Heartbeat(conn, ch_write_msg, ch_err_code)

	b := []byte("['',2,'']")

	for {
		select {
		case msg := <-ch_read_msg:

			var sb []interface{}

			if err := json.Unmarshal(msg, &sb); nil != err {
				log.Fatal(err)
			}

			fmt.Println("data:", sb)

			ch_write_msg <- b

		case code := <-ch_err_code:

			switch code {
			case websocket.CloseNoStatusReceived:
				fmt.Println("CloseNoStatusReceived:", code)
			case websocket.CloseAbnormalClosure:
				fmt.Println("CloseAbnormalClosure:", code)
			case websocket.CloseMessageTooBig:
				fmt.Println("CloseMessageTooBig:", code)
			default:
				fmt.Println("code:", code)
			}

		}
	}
}

func runTcpCli(token string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", *SERVER_ADDR)

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
		Addr:     *REDIS_ADDR,
		Password: *REDIS_PWD,
		DB:       1,
	})

	defer client.Close()

	// pong, err := client.Ping().Result()
	// if nil != err {
	// 	log.Fatal(err)
	// }

	// log.Println(pong)

	uuid := strings.Replace(UUID.Rand().Hex(), "-", "", -1)

	_token, err := client.EvalSha(*REDIS_SHA_AUTH, []string{"1", "1", "backend_1", uuid}, 5, 68, "BACK").Result()
	if nil != err {
		log.Fatal(err)
	}

	token, _ := _token.(string)
	return token
}

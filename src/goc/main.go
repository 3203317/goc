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
	REDIS_PWD      = flag.String("REDIS_PWD", "123456", "Redis密码")
	REDIS_SHA_AUTH = flag.String("REDIS_SHA_AUTH", "a0ad12f31d7de75a5153bdff954caf5bc15b9501", "Redis授权码")
)

var (
	ch_read_msg  = make(chan []byte)
	ch_write_msg = make(chan []byte)
	ch_err       = make(chan error)
)

var (
	ws_url = url.URL{Scheme: "ws", Host: *SERVER_ADDR, Path: "/"}
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

	runWsCli(token)
}

func runWsCli(token string) {
	conn, _, err := websocket.DefaultDialer.Dial(ws_url.String(), nil)

	if nil != err {
		log.Fatal(err)
	}

	defer conn.Close()

	conn.EnableWriteCompression(true)

	go msg.OnMessage(conn, ch_read_msg, ch_err)
	go msg.Process(ch_read_msg, ch_write_msg, ch_err)

	msg.Login(conn, token)
	msg.Heartbeat(conn, ch_write_msg, ch_err)
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

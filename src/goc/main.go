package main

import (
	"config"
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	UUID "github.com/snluu/uuid"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	_ "reflect"
	"strings"
	"time"
)

var (
	SERVER_ADDR    = flag.String("SERVER_ADDR", "47.104.99.102:9988", "前置机地址")
	REDIS_ADDR     = flag.String("REDIS_ADDR", "47.104.99.102:6379", "Redis地址")
	REDIS_PWD      = flag.String("REDIS_PWD", "123456", "Redis密码")
	REDIS_SHA_AUTH = flag.String("REDIS_SHA_AUTH", "", "Redis授权码")

	CLIENT_ID = flag.String("CLIENT_ID", "", "客户端ID")
)

var (
	ch_read_msg  = make(chan []byte)
	ch_write_msg = make(chan []byte)
	ch_err       = make(chan error)
	ch_status    = make(chan Status)
)

var (
	ws_url = url.URL{Scheme: "ws", Host: *SERVER_ADDR, Path: "/"}
)

func def(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

type Status struct {
	code int
	x    x
	msg  string
	data interface{}
}

type x struct {
	x string
}

func main() {
	flag.Parse()

	config.LoadCfg()
	log.Println("Hello, GO!")

	// runHttpServ()

	go start()

	mLoop()
}

func mLoop() {
	for {
		select {
		case status := <-ch_status:

			switch status.code {
			case 0:
			case 1:
				log.Println(1)
				go getToken()
			case 2:
				fmt.Println(2, status.data)
			}

		}
	}
}

func runHttpServ() {
	http.HandleFunc("/", def)
	err := http.ListenAndServe(":80", nil)
	if nil != err {
		log.Fatal(err)
	}
}

func getToken() {
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

	_token, err := client.EvalSha(*REDIS_SHA_AUTH, []string{"1", "1", *CLIENT_ID, uuid}, 5, 68, "BACK").Result()
	if nil != err {
		log.Fatal(err)
	}

	token, _ := _token.(string)

	// fmt.Println(token, reflect.TypeOf(token))
	ch_status <- Status{code: 2, data: token}
}

func start() {
	log.Println("start")

	ticker := time.NewTicker(time.Millisecond * 100)

	defer func() {
		ticker.Stop()
	}()

	select {
	case <-ticker.C:
		ch_status <- Status{code: 1}
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

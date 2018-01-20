package main

import (
	"config"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/gorilla/websocket"
	UUID "github.com/snluu/uuid"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
)

func def(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func main() {
	config.LoadCfg()
	fmt.Println("Hello, GO!")

	token := connRedis()

	fmt.Println(token, reflect.TypeOf(token))

	// http.HandleFunc("/", def)
	// err := http.ListenAndServe(":80", nil)
	// if nil != err {
	// 	log.Fatal(err)
	// }
}

func connRedis() string {
	client := redis.NewClient(&redis.Options{
		Addr:     "47.104.99.102:6379",
		Password: "shuoleniyebudong",
		DB:       1,
	})

	defer client.Close()

	// pong, err := client.Ping().Result()
	// if nil != err {
	// 	log.Fatal(err)
	// }

	// log.Println(pong)

	uuid := strings.Replace(UUID.Rand().Hex(), "-", "", -1)

	_token, err := client.EvalSha("a0ad12f31d7de75a5153bdff954caf5bc15b9501", []string{"1", "1", "backend_1", uuid}, 5, 68, "BACK").Result()
	if nil != err {
		log.Fatal(err)
	}

	token, _ := _token.(string)
	return token
}

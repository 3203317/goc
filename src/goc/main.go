package main

import (
	"config"
	"fmt"
	"github.com/go-redis/redis"
	UUID "github.com/snluu/uuid"
	_ "golang.org/x/net/websocket"
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

	client := redis.NewClient(&redis.Options{
		Addr:     "47.104.99.102:6379",
		Password: "123456",
		DB:       1,
	})

	// pong, err := client.Ping().Result()
	// if nil != err {
	// 	log.Fatal(err)
	// }

	// log.Println(pong)

	token := strings.Replace(UUID.Rand().Hex(), "-", "", -1)

	_token, err := client.EvalSha("a0ad12f31d7de75a5153bdff954caf5bc15b9501", []string{"1", "1", "backend_1", token}, 5, 68, "BACK").Result()
	if nil != err {
		log.Fatal(err)
	}

	log.Println(_token, reflect.TypeOf(_token))

	// http.HandleFunc("/", def)
	// err := http.ListenAndServe(":80", nil)
	// if nil != err {
	// 	log.Fatal(err)
	// }
}

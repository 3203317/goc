package main

import (
	"config"
	"fmt"
	"github.com/go-redis/redis"
	_ "golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
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

	vals, err := client.EvalSha("a0ad12f31d7de75a5153bdff954caf5bc15b9501", []string{"1", "1", "backend_1", "123456"}, 5, 68, "BACK").Result()
	if nil != err {
		log.Fatal(err)
	}

	log.Println(vals)

	// http.HandleFunc("/", def)
	// err := http.ListenAndServe(":80", nil)
	// if nil != err {
	// 	log.Fatal(err)
	// }
}

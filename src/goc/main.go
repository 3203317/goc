package main

import (
	"config"
	"fmt"
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

	http.HandleFunc("/", def)
	err := http.ListenAndServe(":80", nil)
	if nil != err {
		log.Fatal(err)
	}

}

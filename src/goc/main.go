package main

import (
	"config"
	"fmt"
)

func main() {
	config.LoadCfg()
	fmt.Println("Hello, GO!")
}

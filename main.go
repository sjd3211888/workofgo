package main

import (
	"fmt"
	sccliveroom "golearn/websocket"

	"github.com/BurntSushi/toml"
)

func main() {
	//t := make(chan int, 1)

	var conf map[string]map[string]string
	if _, err := toml.DecodeFile("./sccconfig.toml", &conf); err != nil {
		// handle error
	}
	//Host := conf["sccwork"]["Host"]
	fmt.Println("Hostxxxxxxxxxxxxxx", conf)

	go sccliveroom.StartWebsocket()

	select {}

}

package main

import (
	"WebsocketChatroom/global"
	"WebsocketChatroom/server"
	"fmt"
	"log"
	"net/http"
)

var (
	addr   = ":2022"
	banner = `
    ____              _____
   |    |    |   /\     |
   |    |____|  /  \    | 
   |    |    | /----\   |
   |____|    |/      \  |

ChatRoom，start on：%s
`
)

func init() {
	global.Init()
}

func main() {
	fmt.Printf(banner+"\n", addr)
	server.RegisterHandle()
	log.Fatal(http.ListenAndServe(addr, nil))
}

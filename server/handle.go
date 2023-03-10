package server

import (
	"WebsocketChatroom/logic"
	"net/http"
)

func RegisterHandle() {
	go logic.Broadcaster.Start()

	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc)
	http.HandleFunc("/user_list", userListHandleFunc)

}

package server

import (
	"WebsocketChatroom/logic"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func WebSocketHandleFunc(w http.ResponseWriter, req *http.Request) {
	// 如果 Origin 域与主机不同，Accept 将拒绝握手，除非设置了 InsecureSkipVerify 选项（通过第三个参数 AcceptOptions 设置）。
	conn, err := websocket.Accept(w, req, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		log.Println("websocket accept error: ", err)
		return
	}

	// 1. 新用户进来，构建该用户的实例
	token := req.FormValue("token")
	nickname := req.FormValue("nickname")
	if l := len(nickname); l < 2 || l > 20 {
		log.Println("nickname illegal: ", nickname)
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("非法昵称，昵称长度：4-20"))
		conn.Close(websocket.StatusUnsupportedData, "nickname illegal!")
		return
	}

	if !logic.Broadcaster.CanEnterRoom(nickname) {
		log.Println("昵称已存在: ", nickname)
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("该昵称已经存在!"))
		conn.Close(websocket.StatusUnsupportedData, "nickname exists!")
		return
	}

	userHasToken := logic.NewUser(conn, token, nickname, req.RemoteAddr)

	//用户退出后,调用CloseMessageChannel,通道关闭,可使这个协程退出
	go userHasToken.SendMessage(req.Context())

	userHasToken.SendToUserMessageList(logic.NewWelcomeMessage(userHasToken))

	tmpUser := *userHasToken
	user := &tmpUser
	user.Token = ""

	msg := logic.NewUserEnterMessage(user)
	logic.Broadcaster.Broadcast(msg)

	logic.Broadcaster.UserEntering(user)
	log.Println("user: ", nickname, " joins chat")

	err = user.ReceiveMessage(req.Context())

	logic.Broadcaster.UserLeaving(user)
	msg = logic.NewUserLeaveMessage(user)
	logic.Broadcaster.Broadcast(msg)
	log.Println("user: ", nickname, " leaves chat")

	if err == nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Println("read from client error: ", err)
		conn.Close(websocket.StatusInternalError, "Read from client error")
	}

}

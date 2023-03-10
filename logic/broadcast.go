package logic

import (
	"WebsocketChatroom/global"
	"log"
)

var Broadcaster = &broadcaster{
	users: make(map[string]*User),

	enteringChannel: make(chan *User),
	leavingChannel:  make(chan *User),
	messageChannel:  make(chan *Message, global.MessageQueueLen),

	checkUserChannel:      make(chan string),
	checkUserCanInChannel: make(chan bool),

	requestUsersChannel: make(chan struct{}),
	usersChannel:        make(chan []*User),
}

type broadcaster struct {
	users map[string]*User

	enteringChannel chan *User
	leavingChannel  chan *User
	messageChannel  chan *Message

	// 获取用户列表
	requestUsersChannel chan struct{}
	usersChannel        chan []*User

	checkUserChannel      chan string
	checkUserCanInChannel chan bool
}

func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.checkUserChannel <- nickname

	return <-b.checkUserCanInChannel
}

func (b *broadcaster) UserEntering(u *User) {
	b.enteringChannel <- u
}

func (b *broadcaster) UserLeaving(u *User) {
	b.leavingChannel <- u
}

func (b *broadcaster) Broadcast(msg *Message) {
	if len(b.messageChannel) >= global.MessageQueueLen {
		log.Println("broadcast queue 满了")
	}
	b.messageChannel <- msg
}

func (b *broadcaster) Start() {
	for {
		select {
		case user := <-b.enteringChannel:
			b.users[user.NickName] = user

			getOfflineProcessor().Send(user)

		case user := <-b.leavingChannel:
			delete(b.users, user.NickName)

			user.CloseMessageChannel()

		case msg := <-b.messageChannel:
			for _, user := range b.users {
				if user.UID == msg.User.UID {
					continue
				}

				user.SendToUserMessageList(msg)
			}
			getOfflineProcessor().Save(msg)
		case nickname := <-b.checkUserChannel:
			if _, ok := b.users[nickname]; ok {
				b.checkUserCanInChannel <- false
			} else {
				b.checkUserCanInChannel <- true
			}
		case <-b.requestUsersChannel:
			userList := make([]*User, 0, len(b.users))
			for _, user := range b.users {
				userList = append(userList, user)
			}

			b.usersChannel <- userList
		}
	}
}

func (b *broadcaster) GetUserList() []*User {
	b.requestUsersChannel <- struct{}{}
	return <-b.usersChannel
}

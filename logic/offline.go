package logic

import (
	"container/ring"
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

type offlineProcessor struct {
	n int

	// 保存所有用户最近的n条消息
	recentRing *ring.Ring

	// 保存某个用户离线消息（一样 n 条）
	userRing map[string]*ring.Ring
}

var OfflineProcessor *offlineProcessor
var once sync.Once

func getOfflineProcessor() *offlineProcessor {
	once.Do(func() {
		OfflineProcessor = newOfflineProcessor()
	})

	return OfflineProcessor
}

func newOfflineProcessor() *offlineProcessor {
	n := viper.GetInt("offline-num")

	return &offlineProcessor{
		n:          n,
		recentRing: ring.New(n),
		userRing:   make(map[string]*ring.Ring),
	}
}

func (o *offlineProcessor) Save(msg *Message) {
	if msg.Type != MsgTypeNormal {
		return
	}

	fmt.Println(o.recentRing)

	o.recentRing.Value = msg
	o.recentRing = o.recentRing.Next()

	for _, nickname := range msg.Ats {
		nickname = nickname[1:]
		var (
			r  *ring.Ring
			ok bool
		)
		if r, ok = o.userRing[nickname]; !ok {
			r = ring.New(o.n)
		}
		r.Value = msg
		o.userRing[nickname] = r.Next()
	}
}

func (o *offlineProcessor) Send(user *User) {
	fmt.Println(o.recentRing)

	o.recentRing.Do(func(a any) {
		if a != nil {
			user.SendToUserMessageList(a.(*Message))
		}
	})

	if user.isNew {
		return
	}

	if r, ok := o.userRing[user.NickName]; ok {
		r.Do(func(a any) {
			if a != nil {
				user.SendToUserMessageList(a.(*Message))
			}
		})

		delete(o.userRing, user.NickName)
	}
}

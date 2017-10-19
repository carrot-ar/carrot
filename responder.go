package buddy

import (
	"fmt"
	"log"
)

type Responder struct {
	sessions SessionStore

	//inbound messages from the clients
	Broadcast chan []byte
}

func NewResponder() *Responder {
	return &Responder{
		sessions:  NewDefaultSessionManager(),
		Broadcast: make(chan []byte, broadcastChannelSize),
	}
}

func (res *Responder) broadcastAll(message []byte) {
	expiredSessionCount := 0
	closedClientCount := 0
	refreshedClientCount := 0
	messagesSent := 0
	res.sessions.Range(func(key, value interface{}) bool {
		ctx := value.(*Session)

		if ctx.SessionExpired() {
			expiredSessionCount++
			res.sessions.Delete(ctx.Token)
			return true
		} else if !ctx.Client.open {
			closedClientCount++
			return true
		}

		ctx.expireTime = refreshExpiryTime()
		refreshedClientCount++

		select {
		case ctx.Client.send <- message:
			messagesSent++
			return true
		}
	})
	log.Printf("server: broadcast sent %v, refresh %v, closed %v, expired %v",
		messagesSent,
		refreshedClientCount,
		closedClientCount,
		expiredSessionCount)
}

func (res *Responder) Run() {
	for {
		select {
		case message := <-res.Broadcast:
			res.broadcastAll(message)
			fmt.Println(string(message))
		}
	}
}

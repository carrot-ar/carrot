package carrot

import (
	"fmt"
	"log"
)

type Broadcaster struct {
	sessions SessionStore

	//inbound messages from the clients
	broadcast chan []byte
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		sessions:  NewDefaultSessionManager(),
		broadcast: make(chan []byte, broadcastChannelSize),
	}
}

func (br *Broadcaster) broadcastAll(message []byte) {
	expiredSessionCount := 0
	closedClientCount := 0
	refreshedClientCount := 0
	messagesSent := 0
	br.sessions.Range(func(key, value interface{}) bool {
		ctx := value.(*Session)

		if ctx.SessionExpired() {
			expiredSessionCount++
			br.sessions.Delete(ctx.Token)
			return true
		} else if !ctx.Client.Open() {
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

func (br *Broadcaster) Run() {
	for {
		select {
		case message := <-br.broadcast:
			br.broadcastAll(message)
			fmt.Println(string(message))
		}
	}
}

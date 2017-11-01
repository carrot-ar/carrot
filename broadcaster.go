package carrot

import (
	log "github.com/sirupsen/logrus"
)

type Broadcaster struct {
	sessions SessionStore

	//inbound messages from the controllers
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
	log.WithFields(log.Fields{
		"sent":      messagesSent,
		"refreshed": refreshedClientCount,
		"closed":    closedClientCount,
		"expired":   expiredSessionCount,
	}).Debug("broadcast sent")
}

func (br *Broadcaster) Run() {
	for {
		select {
		case message := <-br.broadcast:
			br.broadcastAll(message)
		}
	}
}

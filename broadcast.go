package carrot

import (
//nothing yet
	log "github.com/sirupsen/logrus"
)

type Broadcast struct {
	broadcaster Broadcaster
}

func NewBroadcast(broadcaster Broadcaster) *Broadcast {
	return &Broadcast{
		broadcaster: broadcaster,
	}
}

func (b *Broadcast) Send(message []byte) {
	log.WithField("module", "broadcast").Debug("sending message")
	b.broadcaster.broadcast <- message
}

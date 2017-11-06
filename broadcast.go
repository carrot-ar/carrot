package carrot

import (
	log "github.com/sirupsen/logrus"
)

type Broadcast struct {
	broadcaster Broadcaster
	logger      *log.Entry
}

func NewBroadcast(broadcaster Broadcaster) *Broadcast {
	return &Broadcast{
		broadcaster: broadcaster,
		logger:      log.WithField("module", "broadcast"),
	}
}

func (b *Broadcast) Send(message []byte) {
	b.logger.Debug("sending message")
	b.broadcaster.broadcast <- message
}

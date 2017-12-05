package carrot

import (
	log "github.com/sirupsen/logrus"
)


// Acts as a proxy to the Broadcaster. Messages placed in the Broadcast
// method are broadcasted to the recipients indicated.
type Broadcast struct {
	broadcaster Broadcaster
	logger      *log.Entry
}

// Creates a new instance of the Broadcaster
func NewBroadcast(broadcaster Broadcaster) *Broadcast {
	return &Broadcast{
		broadcaster: broadcaster,
		logger:      log.WithField("module", "broadcast"),
	}
}

// Given a message and a set of sessions, it forwards the message to the
// broadcaster.
func (b *Broadcast) Broadcast(message []byte, sessions ...string) {
	b.logger.Debug("sending message")

	b.broadcaster.broadcast <- OutboundMessage(message, sessions)
}

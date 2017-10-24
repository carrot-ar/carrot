package carrot

import (
	//nothing yet
)

type Broadcast struct {
	broadcaster	*Broadcaster
}

func NewBroadcast(broadcaster *Broadcaster) *Broadcast {
	return &Broadcast {
		broadcaster:	broadcaster,
	}
}

func (b *Broadcast) Send(message []byte) {
	b.broadcaster.broadcast <- message
}
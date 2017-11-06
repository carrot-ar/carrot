package carrot

import (
	log "github.com/sirupsen/logrus"
	"math"
)

var broadcastChannelSize = config.Server.BroadcastChannelSize

type OutboundMessage struct {
	message []byte
}

// manage broadcast groups with the broadcaster
type Broadcaster struct {
	sessions   SessionStore
	clientPool *ClientPool
	//inbound messages from the clients
	broadcast chan []byte
}

func NewBroadcaster(pool *ClientPool) Broadcaster {
	return Broadcaster{
		sessions:   NewDefaultSessionManager(),
		broadcast:  make(chan []byte, broadcastChannelSize),
		clientPool: pool,
	}
}

// use functions are arguments for broadcasting to the correct groups
func (br *Broadcaster) broadcastAll(message []byte) {
	outboundMessage := OutboundMessage{
		// set the criteria function here
		message: message,
	}
	br.clientPool.Send(&outboundMessage)
}

func (br *Broadcaster) Run() {
	for {
		if len(br.broadcast) > int(math.Floor(float64(broadcastChannelSize)*0.90)) {
			log.WithFields(log.Fields{
				"size":   len(br.broadcast),
				"module": "broadcaster"}).Warn("input channel is at or above 90% capacity!")
		}

		if len(br.broadcast) == maxNumDispatcherIncomingRequests {
			log.WithFields(log.Fields{
				"size":   len(br.broadcast),
				"module": "broadcaster"}).Error("input channel is full!")
		}

		select {
		case message := <-br.broadcast:
			log.WithField("module", "broadcaster").Debug("broadcast all")
			br.broadcastAll(message)
		}
	}
}

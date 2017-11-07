package carrot

import (
	log "github.com/sirupsen/logrus"
	"math"
)

const broadcastChannelSize = 4096

type OutboundMessage struct {
	message []byte
}

// manage broadcast groups with the broadcaster
type Broadcaster struct {
	sessions   SessionStore
	clientPool *ClientPool
	//inbound messages from the clients
	broadcast chan []byte
	logger    *log.Entry
}

func NewBroadcaster(pool *ClientPool) Broadcaster {
	return Broadcaster{
		sessions:   NewDefaultSessionManager(),
		broadcast:  make(chan []byte, broadcastChannelSize),
		clientPool: pool,
		logger:     log.WithField("module", "broadcaster"),
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
		if len(br.broadcast) > int(math.Floor(broadcastChannelSize*0.90)) {
			br.logger.WithField("buf_size", len(br.broadcast)).Warn("input channel is at or above 90% capacity!")
		}

		if len(br.broadcast) == maxNumDispatcherIncomingRequests {
			br.logger.WithField("buf_size", len(br.broadcast)).Error("input channel is full!")
		}

		select {
		case message := <-br.broadcast:
			log.WithField("module", "broadcaster").Debug("broadcast all")
			br.broadcastAll(message)
		}
	}
}

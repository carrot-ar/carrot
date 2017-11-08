package carrot

import (
	log "github.com/sirupsen/logrus"
	"math"
)

const broadcastChannelSize = 4096
const broadcastChannelWarningTrigger = 0.9

type OutboundMessage struct {
	message []byte
}

// manage broadcast groups with the broadcaster
type Broadcaster struct {
	sessions SessionStore
	clients  *Clients
	//inbound messages from the clients
	broadcast chan []byte
	logger    *log.Entry
}

func NewBroadcaster(pool *Clients) Broadcaster {
	return Broadcaster{
		sessions:  NewDefaultSessionManager(),
		broadcast: make(chan []byte, broadcastChannelSize),
		clients:   pool,
		logger:    log.WithField("module", "broadcaster"),
	}
}

func (br *Broadcaster) logBufferRedZone() {
	// check for buffer warning
	if len(br.broadcast) > int(math.Floor(broadcastChannelSize*broadcastChannelWarningTrigger)) {
		br.logger.WithField("buf_size", len(br.broadcast)).Warn("input channel is at or above 90% capacity!")
	}
}

func (br *Broadcaster) logBufferFull() {
	// check for buffer full
	if len(br.broadcast) == broadcastChannelSize {
		br.logger.WithField("buf_size", len(br.broadcast)).Error("input channel is full!")
	}
}

func (br *Broadcaster) Run() {
	for {

		br.logBufferRedZone()
		br.logBufferFull()

		select {
		case message := <-br.broadcast:
			for i, client := range br.clients.clients {
				if client.Valid() {

					client.logBufferRedZone()
					client.logBufferFull()

					/*
						TODO: handle full buffers better
						if client.Full() {
							// This can be used to experiment how to handle only writing on our own conditions
							// such as when the buffer size falls below a certain threshold. We can also consider
							// throttling *only* when we reach the red zone or a yellow zone.
						}
					*/

					// **Maintenance Operations**
					// see if the session is expired, if so delete the session.
					// if the client isn't open, then take the client out of the broadcast loop
					// If the client is open, send them the message
					if client.Expired() {
						br.clients.sessions.Delete(client.session.Token)
					} else if !client.Open() { // regardless if the session is expired, see if the client is open or closed.
						br.clients.Release(i)
					}

					// sending operation
					client.session.expireTime = refreshExpiryTime()
					client.send <- message

				} else {
					br.clients.logger.WithFields(log.Fields{
						"i": i,
					}).Debug("nil channel hit!")
				}
			}
		}
	}
}

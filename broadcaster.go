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
			// TODO: Figure out the logic for running a criteria
			// function and only broadcasting to a subset of clients
			for i, client := range br.clients.clients {
				if client.Valid() {

					client.logBufferRedZone()
					client.logBufferFull()

					if client.Full() {
						// This can be used to experiment how to handle only writing on our own conditions
						// such as when the buffer size falls below a certain threshold. We can also consider
						// throttling *only* when we reach the red zone or a yellow zone.
					}

					// see if the session is expired, if so delete the session
					// regardless if the session is expired, see if the client is open or closed.
					// if we are closed, then take the client out of the broadcast loop
					// If the client is open, send them the message
					if client.Expired() {
						br.clients.sessions.Delete(client.session.Token)
					} else if !client.Open() {
						// add the value back to the free list
						// cleanup that slot in the client list
						br.clients.free <- i
						log.WithField("size", len(br.clients.free)).Debugf("adding %v to free list", i)

						br.clients.clients[i] = nil
					} else {
						client.session.expireTime = refreshExpiryTime()
						client.send <- message
					}
				} else {
					br.clients.logger.WithFields(log.Fields{
						"i": i,
					}).Debug("nil channel hit!")
				}
			}
		}
	}
}

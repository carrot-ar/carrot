package carrot

import (
	log "github.com/sirupsen/logrus"
	"math"
)

const broadcastChannelSize = 4096
const broadcastChannelWarningTrigger = 0.9

type OutMessage struct {
	message  []byte
	sessions []string
}

func OutboundMessage(message []byte, sessions []string) OutMessage {
	return OutMessage{
		message:  message,
		sessions: sessions,
	}
}

// manage broadcast groups with the broadcaster
type Broadcaster struct {
	sessions SessionStore
	clients  *Clients
	//inbound messages from the clients
	broadcast chan OutMessage
	logger    *log.Entry
}

func NewBroadcaster(pool *Clients) Broadcaster {
	return Broadcaster{
		sessions:  NewDefaultSessionManager(),
		broadcast: make(chan OutMessage, broadcastChannelSize),
		clients:   pool,
		logger:    log.WithField("module", "broadcaster"),
	}
}

func (br *Broadcaster) checkBufferRedZone() bool {
	// check for buffer warning
	if len(br.broadcast) > int(math.Floor(broadcastChannelSize*broadcastChannelWarningTrigger)) {
		br.logger.WithField("buf_size", len(br.broadcast)).Warn("input channel is at or above 90% capacity!")
		return true

	}

	return false
}

func (br *Broadcaster) checkBufferFull() bool {
	// check for buffer full
	if len(br.broadcast) == broadcastChannelSize {
		br.logger.WithField("buf_size", len(br.broadcast)).Error("input channel is full!")
		return true
	}

	return false
}

func (br *Broadcaster) Run() {
	for {

		br.checkBufferRedZone()
		br.checkBufferFull()

		select {
		case message := <-br.broadcast:
			for i, client := range br.clients.clients {
				if client.Valid() && client.IsRecipient(message.sessions) {

					client.checkBufferRedZone()
					client.checkBufferFull()

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
					client.send <- message.message

				} else {
					//br.clients.logger.WithFields(log.Fields{
					//	"i": i,
					//}).Debug("nil channel hit!")
				}
			}
		}
	}
}

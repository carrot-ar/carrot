package carrot

import (
	"fmt"
	"github.com/DataDog/datadog-go/statsd"
	log "github.com/sirupsen/logrus"
	"math"
)

const broadcastChannelSize = 65536
const broadcastChannelWarningTrigger = 0.9

type OutMessage struct {
	message  []byte
	sessions []string
}

// OutboundMessage initializes a new instance of the OutMessage struct.
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
	statsd    *statsd.Client
}

// NewBroadcaster initializes a new instance of the Broadcaster struct.
func NewBroadcaster(pool *Clients) Broadcaster {
	logger := log.WithField("module", "broadcaster")
	c, err := statsd.New("127.0.0.1:8125")
	if err != nil {
		logger.Error(err)
	}

	return Broadcaster{
		sessions:  NewDefaultSessionManager(),
		broadcast: make(chan OutMessage, broadcastChannelSize),
		clients:   pool,
		logger:    logger,
		statsd:    c,
	}
}

// checkBufferRedZone returns whether the broadcast channel buffer is nearly full.
func (br *Broadcaster) checkBufferRedZone() bool {
	// check for buffer warning
	if len(br.broadcast) > int(math.Floor(broadcastChannelSize*broadcastChannelWarningTrigger)) {
		br.logger.WithField("buf_size", len(br.broadcast)).Warn("input channel is at or above 90% capacity!")
		return true

	}

	return false
}

// checkBufferFull returns whether the broadcast channel buffer is actually full.
func (br *Broadcaster) checkBufferFull() bool {
	// check for buffer full
	if len(br.broadcast) == broadcastChannelSize {
		br.logger.WithField("buf_size", len(br.broadcast)).Error("input channel is full!")
		return true
	}

	return false
}

//Run sends buffered responses to devices, deletes expired device connections, and logs information.
func (br *Broadcaster) Run() {
	for {
		br.statsd.Gauge("carrot.broadcaster.outbound.buffer_size", float64(len(br.broadcast)), nil, 100)

		br.checkBufferRedZone()
		br.checkBufferFull()

		select {
		case message := <-br.broadcast:
			for i, client := range br.clients.clients {
				if client.Valid() && client.IsRecipient(message.sessions) {
					br.clients.logger.WithFields(log.Fields{
						"i": i,
					}).Debug("valid channel hit!")

					statsdName := fmt.Sprintf("carrot.client.%v.buffer_size", i)
					client.statsd.Gauge(statsdName, float64(len(client.send)), nil, 100)

					client.checkBufferRedZone()
					client.checkBufferFull()

					// **Maintenance Operations**
					// see if the session is expired, if so delete the session.
					// if the client isn't open, then take the client out of the broadcast loop
					// If the client is open, send them the message
					if client.Expired() {
						br.clients.sessions.Delete(client.session.Token)
					} else if !client.Open() { // regardless if the session is expired, see if the client is open or closed.
						br.clients.Release(i)
						continue
					}

					// sending operation
					client.session.expireTime = refreshExpiryTime()
					client.logger.Infof("client is %v", client.Open())
					client.send <- message.message

				} else {
					br.clients.logger.WithFields(log.Fields{
						"i": i,
					}).Debug("nil channel hit!")
				}
			}
		}
	}
}

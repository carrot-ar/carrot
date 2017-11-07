package carrot

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
)

const maxClients = 16
const maxClientPoolQueueBackup = 128
const maxOutboundMessages = 4096

type ClientPool struct {
	sessions             SessionStore
	clients              []*Client
	free                 chan int
	insertQueue          chan *Client
	outboundMessageQueue chan *OutboundMessage
	logger               *log.Entry
}

func NewClientPool() *ClientPool {
	// setup the free list by filling up a channel of
	// integers from 0 to maxClients
	free := make(chan int, maxClients)
	for i := 0; i < maxClients; i++ {
		free <- i
	}

	return &ClientPool{
		sessions:             NewDefaultSessionManager(),
		clients:              make([]*Client, maxClients, maxClients),
		free:                 free,
		insertQueue:          make(chan *Client, maxClientPoolQueueBackup),
		outboundMessageQueue: make(chan *OutboundMessage, maxOutboundMessages),
		logger:               log.WithField("module", "client_pool"),
	}
}

// Add to the free list in a non blocking fashion
// by adding the client to a queue that will insert new clients
// on the broadcast loop
func (cp *ClientPool) Insert(client *Client) error {
	if len(cp.insertQueue) == maxClientPoolQueueBackup {
		return fmt.Errorf("unable to queue client")
	}

	cp.insert(client)

	return nil
}

// potential data race could develop here. Test with -race
//
// Grab from the insertQueue and set the client. Will be called in the broadcast loop
func (cp *ClientPool) insert(client *Client) error {
	var index int
	index, err := cp.getFreeIndex()
	if err != nil {
		return err
	}

	err = cp.setClient(index, client)
	if err != nil {
		return err
	}

	return nil
}

func (cp *ClientPool) setClient(index int, client *Client) error {
	if cp.clients[index] != nil {
		return fmt.Errorf("index %v contained a client when it should not have", index)
	}

	cp.clients[index] = client
	return nil
}

// Get the a free index from the free list
func (cp *ClientPool) getFreeIndex() (int, error) {
	if len(cp.free) == 0 {
		return -1, fmt.Errorf("client pool is full")
	}

	freeSpot := <-cp.free
	return freeSpot, nil
}

func (cp *ClientPool) Count() int {
	return 1 - len(cp.free)
}

// send a message to the clients
func (cp *ClientPool) Send(message *OutboundMessage) {
	cp.logger.Debug("in the send method")
	cp.outboundMessageQueue <- message
}

// loop and send
func (cp *ClientPool) ListenAndSend() {
	for {

		if len(cp.insertQueue) > int(math.Floor(maxClientPoolQueueBackup*0.90)) {
			cp.logger.WithFields(log.Fields{
				"size":    len(cp.insertQueue),
				"channel": "insert_queue"}).Warn("input channel is at or above 90% capacity!")
		}

		if len(cp.insertQueue) == maxClientPoolQueueBackup {
			cp.logger.WithFields(log.Fields{
				"size":    len(cp.insertQueue),
				"channel": "insert_queue"}).Error("input channel is full!")
		}

		if len(cp.outboundMessageQueue) > int(math.Floor(maxOutboundMessages*0.90)) {
			cp.logger.WithFields(log.Fields{
				"size":    len(cp.outboundMessageQueue),
				"channel": "outbound"}).Warn("input channel is at or above 90% capacity!")
		}

		if len(cp.outboundMessageQueue) == maxOutboundMessages {
			cp.logger.WithFields(log.Fields{
				"size":    len(cp.outboundMessageQueue),
				"channel": "outbound"}).Error("input channel is full!")
		}

		select {
		case newClient := <-cp.insertQueue:
			cp.insert(newClient)
		case message := <-cp.outboundMessageQueue:
			// TODO: Figure out the logic for running a criteria
			// function and only broadcasting to a subset of clients
			for i, client := range cp.clients {
				if client != nil {

					log.WithFields(log.Fields{
						"i":     i,
						"open?": client.Open(),
					}).Debug("open channel")

					if len(client.send) > int(math.Floor(sendMsgBufferSize*0.90)) {
						cp.logger.WithFields(log.Fields{
							"i":       i,
							"open?":   client.Open(),
							"size":    len(client.send),
							"channel": "send"}).Warn("input channel is at or above 90% capacity!")
					}

					if len(client.send) == sendMsgBufferSize {
						cp.logger.WithFields(log.Fields{
							"i":       i,
							"open?":   client.Open(),
							"size":    len(client.send),
							"channel": "send"}).Error("input channel is full!")
					}

					// see if the session is expired, if so delete the session
					// regardless if the session is expired, see if the client is open or closed.
					// if we are closed, then take the client out of the broadcast loop
					// If the client is open, send them the message
					if client.Expired() {
						cp.sessions.Delete(client.session.Token)
					} else if !client.Open() {
						// add the value back to the free list
						// cleanup that slot in the client list
						cp.free <- i
						log.WithField("size", len(cp.free)).Infof("adding %v to free list", i)

						// the client variable seems to be a copy of the value of cp.clients[i]? we want
						// to modify the pointer. Strange behavior
						// TODO: investigate
						cp.clients[i] = nil
					} else {
						client.session.expireTime = refreshExpiryTime()
						client.send <- message.message
					}

				} else {
					cp.logger.WithFields(log.Fields{
						"i": i,
					}).Debug("nil channel hit!")
				}
			}

		default:
			//log.Debug("DEFAULT ENTRY HIT")
		}
	}
}

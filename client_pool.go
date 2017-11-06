package carrot

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
)

var (
	clientPoolConfig         = config.ClientPool
	maxClients               = clientPoolConfig.MaxClients
	maxClientPoolQueueBackup = clientPoolConfig.MaxClientPoolQueueBackup
	maxOutboundMessages      = clientPoolConfig.MaxOutboundMessages
)

type ClientPool struct {
	sessions             SessionStore
	clients              []*Client
	free                 chan int
	insertQueue          chan *Client
	outboundMessageQueue chan *OutboundMessage
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
	log.Debug("in the send method")
	cp.outboundMessageQueue <- message
}

// loop and send
func (cp *ClientPool) ListenAndSend() {
	for {

		if len(cp.insertQueue) > int(math.Floor(float64(maxClientPoolQueueBackup)*0.90)) {
			log.WithFields(log.Fields{
				"size":    len(cp.insertQueue),
				"module":  "client_pool",
				"channel": "insert_queue"}).Warn("input channel is at or above 90% capacity!")
		}

		if len(cp.insertQueue) == maxClientPoolQueueBackup {
			log.WithFields(log.Fields{
				"size":    len(cp.insertQueue),
				"module":  "client_pool",
				"channel": "insert_queue"}).Error("input channel is full!")
		}

		if len(cp.outboundMessageQueue) > int(math.Floor(float64(maxOutboundMessages)*0.90)) {
			log.WithFields(log.Fields{
				"size":    len(cp.outboundMessageQueue),
				"module":  "client_pool",
				"channel": "outbound"}).Warn("input channel is at or above 90% capacity!")
		}

		if len(cp.outboundMessageQueue) == maxOutboundMessages {
			log.WithFields(log.Fields{
				"size":    len(cp.outboundMessageQueue),
				"module":  "client_pool",
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
						"module": "client_pool",
						"i":      i,
						"open?":  client.Open(),
					}).Debug("open channel")

					if len(client.send) > int(math.Floor(float64(sendMsgBufferSize)*0.90)) {
						log.WithFields(log.Fields{
							"size":    len(client.send),
							"module":  "client",
							"channel": "send"}).Warn("input channel is at or above 90% capacity!")
					}

					if len(client.send) == sendMsgBufferSize {
						log.WithFields(log.Fields{
							"size":    len(client.send),
							"module":  "client",
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
					log.WithFields(log.Fields{
						"i":      i,
						"module": "client_pool",
					}).Debug("nil channel hit!")
				}
			}

		default:
			//log.Debug("DEFAULT ENTRY HIT")
		}
	}
}

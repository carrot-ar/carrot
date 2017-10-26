package carrot

import "fmt"

const maxClientPoolQueueBackup = 128

type Pool interface{
	Insert(*Client) error
	Count() int
	Send(*OutboundMessage)
	ListenAndSend()
}

type ClientPool struct {
	sessions SessionStore
	clients []*Client
	free	chan int
	insertQueue chan *Client
	outboundMessageQueue chan *OutboundMessage
}

func NewClientPool() Pool {
	// setup the free list by filling up a channel of
	// integers from 0 to maxClients
	free := make(chan int, maxClients)
	for i := 0; i < maxClients; i++ {
		free <- i
	}

	return &ClientPool{
		sessions: NewDefaultSessionManager(),
		clients: make([]*Client, maxClients, maxClients),
		free: free,
		insertQueue: make(chan *Client, maxClientPoolQueueBackup),
	}
}


// Add to the free list in a non blocking fashion
// by adding the client to a queue that will insert new clients
// on the broadcast loop
func (cp *ClientPool) Insert(client *Client) error {
	if len(cp.insertQueue) == maxClientPoolQueueBackup {
		return fmt.Errorf("unable to queue client")
	}

	cp.insertQueue <- client

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

	err = cp.setClient(index, <-cp.insertQueue)
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
	return <-cp.free, nil
}

func (cp *ClientPool) Count() int {
	return 1 - len(cp.free)
}

// send a message to the clients
func (cp *ClientPool) Send(message *OutboundMessage) {
	cp.outboundMessageQueue <- message
}

// loop and send
func (cp *ClientPool) ListenAndSend() {
	for {
		select {
		case message := <- cp.outboundMessageQueue:
			// TODO: Figure out the logic for running a criteria
			// function and only broadcasting to a subset of clients
			for i, client := range cp.clients {
				if client != nil {

					if client.Expired() {
						cp.sessions.Delete(client.session.Token)
						continue
					} else if !client.Open() {
						// add the value back to the free list
						// cleanup that slot in the client list
						cp.free <- i
						client = nil
						continue
					}

					client.session.expireTime = refreshExpiryTime()
					client.send <- message.message

				}
			}
		case newClient := <- cp.insertQueue:
			cp.insert(newClient)
		}
	}
}

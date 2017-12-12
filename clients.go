package carrot

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
)

const maxClients = 16

type Clients struct {
	sessions SessionStore
	clients  []*Client
	free     chan int
	mutex    *sync.RWMutex
	logger   *log.Entry
}

func NewClientList() (*Clients, error) {
	if maxClients < 1 {
		return nil, errors.New("client list size must be greater than 0")
	}

	// setup the free list by filling up a channel of
	// integers from 0 to maxClients
	free := make(chan int, maxClients)
	for i := 0; i < maxClients; i++ {
		free <- i
	}

	return &Clients{
		sessions: NewDefaultSessionManager(),
		clients:  make([]*Client, maxClients, maxClients),
		free:     free,
		logger:   log.WithField("module", "client_pool"),
		mutex:    &sync.RWMutex{},
	}, nil
}

// Thread-safe add to the client list
func (cp *Clients) Insert(client *Client) error {

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

/*
 * Releases an index in the client list. The index is added back onto the free list
 * then the slot in the client list is set to nil.
 * This more or less a maintenance operation that will occur is a client is not longer open
 */
func (cp *Clients) Release(index int) {
	log.WithField("size", len(cp.free)).Debugf("releasing %v", index)

	cp.free <- index

	cp.mutex.Lock()
	cp.clients[index] = nil
	cp.mutex.Unlock()
}

func (cp *Clients) setClient(index int, client *Client) error {
	if cp.clients[index] != nil {
		return fmt.Errorf("index %v contained a client when it should not have", index)
	}

	cp.mutex.Lock()
	cp.clients[index] = client
	cp.mutex.Unlock()

	return nil
}

// Get the a free index from the free list
func (cp *Clients) getFreeIndex() (int, error) {
	if len(cp.free) == 0 {
		return -1, fmt.Errorf("client pool is full")
	}

	free := <-cp.free
	return free, nil
}

/*
func (cp *Clients) Count() int {
	return 1 - len(cp.free)
}
*/

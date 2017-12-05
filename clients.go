package carrot

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
)

const maxClients = 16

/*
	Clients is the data structure used to keep track of individual connections
	to the Carrot server. The Clients structure maintenance occurs within the Broadcaster
	which handles sending data to individual connections.
 */
type Clients struct {
	// Pointer to the global session store
	sessions SessionStore

	// Client list
	clients  []*Client

	// Free list used for choosing which slot to add a new client to
	free     chan int

	// Mutex for client list access
	mutex    *sync.RWMutex

	logger   *log.Entry
}

/*
	Creates a new Clients structure and builds a free list.
 */
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

// Thread-safe addition of a client to the client list.
func (cp *Clients) Insert(client *Client) error {

	var index int
	index, err := cp.getFreeIndex()
	if err != nil {
		return err
	}

	cp.mutex.Lock()
	err = cp.setClient(index, client)
	cp.mutex.Unlock()
	if err != nil {
		return err
	}

	return nil
}


// Releases an index in the client list. This happens by adding the entry back into the free list,
// then the slot in the client list is set to nil.
func (cp *Clients) Release(index int) {
	log.WithField("size", len(cp.free)).Debugf("releasing %v", index)
	cp.free <- index
	cp.clients[index] = nil
}

func (cp *Clients) setClient(index int, client *Client) error {
	if cp.clients[index] != nil {
		return fmt.Errorf("index %v contained a client when it should not have", index)
	}

	cp.clients[index] = client
	return nil
}

// Get the a free index from the free list
func (cp *Clients) getFreeIndex() (int, error) {
	if len(cp.free) == 0 {
		return -1, fmt.Errorf("client pool is full")
	}

	freeSpot := <-cp.free
	return freeSpot, nil
}

/*
func (cp *Clients) Count() int {
	return 1 - len(cp.free)
}
*/

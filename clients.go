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

	cp.mutex.Lock()
	err = cp.setClient(index, client)
	cp.mutex.Unlock()
	if err != nil {
		return err
	}

	return nil
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

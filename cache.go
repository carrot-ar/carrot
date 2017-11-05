package carrot

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"time"
)

type CachedControllersList struct {
	cachedControllers map[string]*AppController
	lru               *PriorityQueue
}

func NewCachedControllersList() *CachedControllersList {
	return &CachedControllersList{
		cachedControllers: make(map[string]*AppController),
		lru:               NewPriorityQueue(),
	}
}

func (ccl *CachedControllersList) Exists(key string) bool {
	_, ok := ccl.cachedControllers[key]
	if ok {
		return true
	}
	return false
}

func (ccl *CachedControllersList) Get(key string) (*AppController, error) {
	var err error
	cc, ok := ccl.cachedControllers[key]
	if !ok || cc == nil {
		log.WithFields(log.Fields{
			"session_token": key,
		})
		err = errors.New("cached controller does not exist")
	}
	//update priority to reflect recent controller usage
	ccl.lru.UpdatePriority(key, getPriority())
	return cc, err
}

func (ccl *CachedControllersList) Add(key string, ac *AppController) {
	//add to controller map
	ccl.cachedControllers[key] = ac
	//add to LRU to keep track of deletion order
	ccl.lru.Insert(key, getPriority())
	log.WithFields(log.Fields{
		"key":        key,
		"cache_size": ccl.lru.Len(),
	}).Debug("new controller cached")
}

func (ccl *CachedControllersList) DeleteOldest() error {
	//find oldest controller in LRU to identify token and delete LRU record
	key, err := ccl.lru.Pop()
	if err != nil {
		return err
	}
	//use token to delete from controller map
	delete(ccl.cachedControllers, key.(string)) //doesn't return anything
	return nil
}

func (ccl *CachedControllersList) IsEmpty() bool {
	if len(ccl.cachedControllers) == 0 {
		return true
	}
	return false
}

func (ccl *CachedControllersList) Length() int {
	return len(ccl.cachedControllers)
}

func getPriority() float64 {
	return float64(time.Now().Unix())
}

package carrot

import (
	"errors"
	"time"
)

type CachedControllersList struct {
	cachedControllers map[string]*AppController
	lru               *PriorityQueue
}

// NewCachedControllersList initializes a new instance of the CachedControllersList struct.
func NewCachedControllersList() *CachedControllersList {
	return &CachedControllersList{
		cachedControllers: make(map[string]*AppController),
		lru:               NewPriorityQueue(),
	}
}

// Exists checks whether a controller of the given type is already being cached.
func (ccl *CachedControllersList) Exists(key string) bool {
	_, ok := ccl.cachedControllers[key]
	if ok {
		return true
	}
	return false
}

// Get returns a reference to the controller whose type matches the key.
func (ccl *CachedControllersList) Get(key string) (*AppController, error) {
	var err error
	cc, ok := ccl.cachedControllers[key]
	if !ok || cc == nil {
		err = errors.New("cached controller does not exist")
	}

	//update priority to reflect recent controller usage
	ccl.lru.UpdatePriority(key, getPriority())
	return cc, err
}

// Add maintains the existence of a new controller so that it does not have to be created again upon request reference.
func (ccl *CachedControllersList) Add(key string, ac *AppController) {
	//add to controller map
	ccl.cachedControllers[key] = ac
	//add to LRU to keep track of deletion order
	ccl.lru.Insert(key, getPriority())
}

// DeleteOldest frees the memory of the oldest controller maintained by the cache.
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

// IsEmpty determines whether the cache is currently maintaining any controllers.
func (ccl *CachedControllersList) IsEmpty() bool {
	if len(ccl.cachedControllers) == 0 {
		return true
	}
	return false
}

// Length returns the size of the cache: the number of maintained controllers.
func (ccl *CachedControllersList) Length() int {
	return len(ccl.cachedControllers)
}

// getPriority returns the current Unix time in order to record controller age for later deletion.
func getPriority() float64 {
	return float64(time.Now().Unix())
}

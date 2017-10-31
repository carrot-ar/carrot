package carrot

import (
	"fmt"
	"time"
)

type CachedControllersList struct {
	cachedControllers map[string]*AppController
	lru *PriorityQueue
}

func NewCachedControllersList() *CachedControllersList {
	return &CachedControllersList{
		cachedControllers: make(map[string]*AppController),
		lru:	NewPriorityQueue(),
	}
}

func (ccl *CachedControllersList) Exists(key string) bool {
	_, ok := ccl.cachedControllers[key]
	if ok {
		return true
	}
	return false
}

func (ccl *CachedControllersList) Get(key string) *AppController {
	cc, ok := ccl.cachedControllers[key]
	if !ok || cc == nil {
		fmt.Println("cannot return route because it doesn't exist")
		//return nil
	}

	//update priority to reflect recent controller usage
	ccl.lru.UpdatePriority(key, getPriority())

	return cc
}

func (ccl *CachedControllersList) Add(key string, ac *AppController) {
	//add to controller map
	ccl.cachedControllers[key] = ac
	//add to LRU to keep track of deletion order
	ccl.lru.Insert(key, getPriority())
	fmt.Printf("a controller has been added, num of controllers: %v \n", ccl.lru.Len())
}

func (ccl *CachedControllersList) DeleteOldest() {
	//find oldest controller in LRU to identify token and delete LRU record
	key, err := ccl.lru.Pop()
	if (err != nil) {
		fmt.Println("there was a problem removing the oldest element from the LRU")
	}

	//use token to delete from controller map
	delete(ccl.cachedControllers, key.(string)) //doesn't return anything
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
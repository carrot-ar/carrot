package carrot

import (
	log "github.com/sirupsen/logrus"
)

type CachedControllersList struct {
	cachedControllers map[SessionToken]*AppController
}

func NewCachedControllersList() *CachedControllersList {
	return &CachedControllersList{
		cachedControllers: make(map[SessionToken]*AppController),
	}
}

func (ccl *CachedControllersList) Exists(token SessionToken) bool {
	_, ok := ccl.cachedControllers[token]
	if ok {
		return true
	}
	return false
}

func (ccl *CachedControllersList) Get(token SessionToken) *AppController {
	cc, ok := ccl.cachedControllers[token]
	if !ok || cc == nil {
		log.WithFields(log.Fields{
			"session_token" : token,
		}).Error("cached controller does not exist")
	}
	return cc
}

func (ccl *CachedControllersList) Add(token SessionToken, ac *AppController) {
	ccl.cachedControllers[token] = ac
}

func (ccl *CachedControllersList) Delete(token SessionToken) {
	delete(ccl.cachedControllers, token) //doesn't return anything
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

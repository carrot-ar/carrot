package carrot

import (
	log "github.com/sirupsen/logrus"
)

var broadcast *Broadcast
var Environment string

// TODO: refactor so that if a module fails to load, we cause an error
func Run() error {

	// if the environment isn't set, then we can set to debug.
	if Environment != "testing" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.PanicLevel)
	}

	sessions := NewDefaultSessionManager()
	log.Debug("session store initialized")
	clientPool := NewClientPool()
	log.Debug("client pool initialized")
	server := NewServer(clientPool, sessions)
	log.Debug("server initialized")
	dispatcher := NewDispatcher()
	log.Debug("dispatcher initialized")
	broadcaster := NewBroadcaster(clientPool)
	log.Debug("broadcaster initialized")
	broadcast = NewBroadcast(broadcaster)
	log.Debug("set broadcast var")
	go dispatcher.Run()
	log.Debug("dispatcher started")
	go server.Middleware.Run()
	log.Debug("middleware started")
	go server.Run()
	log.Debug("server started")
	go broadcast.broadcaster.clientPool.ListenAndSend()
	log.Debug("client pool broadcaster started")
	log.Debug("beginning to serve")

	if Environment != "testing" {
		server.Serve()
	}

	return nil
}

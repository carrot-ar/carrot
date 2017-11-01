package carrot

import (
	log "github.com/sirupsen/logrus"
)

var Environment string
var broadcaster Broadcaster

// TODO: refactor so that if a module fails to load, we cause an error
func Run() error {

	// if the environment isn't set, then we can set to debug.
	if Environment != "testing" {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.PanicLevel)
	}

	sessions := NewDefaultSessionManager()
	log.Debug("session store initialized")
	clientPool := NewClientPool()
	log.Debug("client pool initialized")
	server := NewServer(clientPool, sessions)
	log.Debug("server initialized")

	// TODO: clean all this up
	broadcaster = NewBroadcaster(clientPool)
	log.Debug("global broadcaster created")
	go broadcaster.clientPool.ListenAndSend()
	log.Debug("global broadcaster running")
	go server.Middleware.Run()
	log.Debug("middleware started")
	go server.Run()
	log.Debug("server started")
	log.Debug("beginning to serve")

	if Environment != "testing" {
		server.Serve()
	}

	return nil
}

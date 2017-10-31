package carrot

import (
	log "github.com/sirupsen/logrus"
)

var Environment string

func Run() error {
	// if the environment isn't set, then we can set to debug.
	if Environment != "testing" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.PanicLevel)
	}

	// TODO: refactor so that if a module fails to load, we cause an error
	sessions := NewDefaultSessionManager()
	server := NewServer(sessions)
	dispatcher := NewDispatcher()
	go dispatcher.Run()
	go server.Middleware.Run()
	go server.Run()
	log.Debug("server started")
	log.Debug("beginning to serve")

	if Environment != "testing" {
		server.Serve()
	}

	return nil
}

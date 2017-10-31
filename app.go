package carrot

import (
	log "github.com/sirupsen/logrus"
)

var Environment string
	
	


// TODO: refactor so that if a module fails to load, we cause an error
func Run() {
  // if the environment isn't set, then we can set to debug.
  
  if Environment != "testing" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.PanicLevel)
	}
  
  sessions := NewDefaultSessionManager()
	log.Debug("session store initialized")
	server := NewServer(sessions)
	log.Debug("server initialized")
	dispatcher := NewDispatcher()
	log.Debug("dispatcher initialized")
	go dispatcher.Run()
	log.Debug("dispatcher started")
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

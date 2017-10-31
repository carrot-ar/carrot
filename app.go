package carrot

import (
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.InfoLevel)
}

func Run() {
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
	server.Serve()
}

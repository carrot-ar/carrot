package carrot

import (
	log "github.com/sirupsen/logrus"
)

const (
	//for dispatcher_test.go
	endpoint1 = "test1"
)

func init() {
	log.SetLevel(log.PanicLevel)

	//for dispatcher_test.go
	Add(endpoint1, TestDispatcherController{}, "Print", false)	
}

package carrot

import (
	log "github.com/sirupsen/logrus"
)

/*
	Response groups are enumerated by the constants in this file.
	To add a new response group,  add it to the list of constants
	at the start of the file.
*/
type Broadcast struct {
	broadcaster Broadcaster
	logger      *log.Entry
}

// NewBroadcast initializes a new instance of the Broadcast struct.
func NewBroadcast(broadcaster Broadcaster) *Broadcast {
	return &Broadcast{
		broadcaster: broadcaster,
		logger:      log.WithField("module", "broadcast"),
	}
}

/*
 // ...string []string

 Broadcast(message, "sessiontoken", "sessiontoken")
 Broadcast(message, sessionTokens)
 Broadcast(message)

 Broadcast(message, req.SessionToken)

 func (c *TestController) EchoToSomePeople(req *carrot.Request, broadcast *carrot.Broadcast) {
    response := carrot.Response()
    response.AddParam("hello": "world")
    response.Build()
 	responseGroup = mysql.find(/*query to find users).SessionTokens
    broadcast.Broadcast(message, responseGroup)
 }
*/

// Broadcast allows controllers to send responses to all or a subset of listening devices.
func (b *Broadcast) Broadcast(message []byte, sessions ...string) {
	b.logger.Debug("sending message")

	b.broadcaster.broadcast <- OutboundMessage(message, sessions)
}

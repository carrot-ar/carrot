package buddy

import (
	//"fmt"
	//"log"
	"time"
)

type Request struct {
	session   *Session
	message   []byte
	startTime time.Time
}

func NewRequest(session *Session, message []byte) *Request {
	return &Request{
		session:   session,
		message:   message,
		startTime: time.Now(),
	}
}

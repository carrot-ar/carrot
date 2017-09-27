package buddy

import (
//"fmt"
//"log"
)

type Request struct {
	session *Session
	message []byte
}

func NewRequest(session *Session, message []byte) *Request {
	return &Request{
		session: session,
		message: message,
	}
}

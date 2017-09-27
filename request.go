package buddy

import (
	//"fmt"
	//"log"
)

type Request struct {
	session *Session
	message string
}

func NewRequest(session *Session, message string) {
	return &Request{
			
	}
}


package buddy

import (
	"testing"
)

func TestBroadcasting(t *testing.T) {
	responder := NewResponder()
	go responder.Run()
	responder.Broadcast <- []byte("This is the responder test broadcasting a message!")
}

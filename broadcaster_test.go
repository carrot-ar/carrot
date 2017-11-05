package carrot

import (
	"testing"
)

func TestBroadcasting(t *testing.T) {
	clientPool := NewClientPool()
	broadcaster := NewBroadcaster(clientPool)
	broadcast := NewBroadcast(broadcaster)
	go broadcast.broadcaster.Run()

	broadcast.Send([]byte("This is the broadcaster test broadcasting a message!"))
}

package carrot

import (
	log "github.com/sirupsen/logrus"
	"math"
	"testing"
)

func TestBroadcasting(t *testing.T) {
	clientPool, _ := NewClientList()
	broadcaster := NewBroadcaster(clientPool)
	broadcast := NewBroadcast(broadcaster)
	go broadcast.broadcaster.Run()

	broadcast.Send([]byte("This is the broadcaster test broadcasting a message!"))
}

func TestBroadcastercheckBufferRedzone(t *testing.T) {
	broadcaster := &Broadcaster{
		broadcast: make(chan []byte, broadcastChannelSize),
		logger:    log.WithField("module", "broadcaster_test"),
	}

	res := broadcaster.checkBufferRedZone()

	if res == true {
		t.Fatalf("buffer was signaling redzone when it shouldn't have been")
	}

	for i := 0; i < int(math.Floor(broadcastChannelSize*0.95)); i++ {
		broadcaster.broadcast <- ([]byte("test message!"))
	}

	res = broadcaster.checkBufferRedZone()

	if res != true {
		t.Fatalf("buffer was not signaling redzone when it should have been ")
	}
}

func TestBroadcastercheckBufferFull(t *testing.T) {
	broadcaster := &Broadcaster{
		broadcast: make(chan []byte, broadcastChannelSize),
		logger:    log.WithField("module", "broadcaster_test"),
	}

	res := broadcaster.checkBufferFull()

	if res == true {
		t.Fatalf("buffer was signaling full when it shouldn't have been")
	}

	for i := 0; i < broadcastChannelSize; i++ {
		broadcaster.broadcast <- ([]byte("test message!"))
	}

	res = broadcaster.checkBufferFull()

	if res != true {
		t.Fatalf("buffer was not full redzone when it should have been ")
	}
}

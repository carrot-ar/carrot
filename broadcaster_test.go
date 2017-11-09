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

	broadcast.Broadcast([]byte("This is the broadcaster test broadcasting a message!"))
}

func TestBroadcastercheckBufferRedzone(t *testing.T) {
	broadcaster := &Broadcaster{
		broadcast: make(chan OutMessage, broadcastChannelSize),
		logger:    log.WithField("module", "broadcaster_test"),
	}

	res := broadcaster.checkBufferRedZone()

	if res == true {
		t.Fatalf("buffer was signaling redzone when it shouldn't have been")
	}

	for i := 0; i < int(math.Floor(broadcastChannelSize*0.95)); i++ {
		broadcaster.broadcast <- OutboundMessage([]byte("test message!"), []string{})
	}

	res = broadcaster.checkBufferRedZone()

	if res != true {
		t.Fatalf("buffer was not signaling redzone when it should have been ")
	}
}

func TestBroadcastercheckBufferFull(t *testing.T) {
	broadcaster := &Broadcaster{
		broadcast: make(chan OutMessage, broadcastChannelSize),
		logger:    log.WithField("module", "broadcaster_test"),
	}

	res := broadcaster.checkBufferFull()

	if res == true {
		t.Fatalf("buffer was signaling full when it shouldn't have been")
	}

	for i := 0; i < broadcastChannelSize; i++ {
		broadcaster.broadcast <- OutboundMessage([]byte("test message!"), []string{})
	}

	res = broadcaster.checkBufferFull()

	if res != true {
		t.Fatalf("buffer was not full redzone when it should have been ")
	}
}

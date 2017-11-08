package carrot

import (
	"sync"
	"testing"
	"time"
)

func sampleClient() *Client {
	sessions := NewDefaultSessionManager()
	_, session, _ := sessions.NewSession()

	client := &Client{
		session:   session,
		send:      make(chan []byte, sendMsgBufferSize),
		sendToken: make(chan SessionToken, sendTokenBufferSize),
		start:     make(chan struct{}),
		openMutex: &sync.RWMutex{},
		open:      true,
	}

	return client
}

func TestClientExpired(t *testing.T) {
	client := sampleClient()

	waitTime := time.Now()

	if client.Expired() == true {
		t.Fatalf("client expired when it shouldn't be")
	}

	client.session.expireTime = waitTime
	client.softClose()

	if client.Expired() == false {
		t.Fatalf("client should be expired but it isn't")
	}
}

func TestClientFull(t *testing.T) {
	client := sampleClient()

	if client.Full() == true {
		t.Fatalf("client send buffer is full when it should be empty size: %v", len(client.sendToken))
	}

	for i := 0; i < sendMsgBufferSize; i++ {
		client.send <- []byte("dummy message")
	}

	if client.Full() == false {
		t.Fatalf("client send buffer is not full when it should be! size: %v", len(client.sendToken))
	}
}

func TestClientOpen(t *testing.T) {
	client := sampleClient()

	client.softOpen()

	if client.Open() == false {
		t.Fatal("client is closed when it should be open")
	}

	client.softClose()
	if client.Open() == true {
		t.Fatal("client is open when it should be closed")
	}
}

func TestClientValid(t *testing.T) {
	client := sampleClient()

	if client.Valid() != true {
		t.Fatal("client is invalid when it should be valid")
	}

	client = nil

	if client.Valid() != false {
		t.Fatal("client is valid when it should not be valid")
	}
}

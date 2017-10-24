package carrot

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"testing"
	"time"
)

func TestDefaultSessionManagerNewSessionGet(t *testing.T) {
	store := NewDefaultSessionManager()

	token, err := store.NewSession()
	if err != nil {
		t.Errorf("Failed to create session")
	}

	_, err = store.Get(token)
	if err != nil {
		t.Errorf("Failed to get session %v", err)
	}
}

func TestContextPersistence(t *testing.T) {
	store := NewDefaultSessionManager()

	token, err := store.NewSession()
	if err != nil {
		t.Errorf("Failed to create session")
	}

	ctx, _ := store.Get(token)

	if ctx == nil {
		t.Error("Session was not received")
	}
}

func TestSessionDelete(t *testing.T) {
	store := NewDefaultSessionManager()

	token, err := store.NewSession()
	if err != nil {
		t.Errorf("Failed to create session")
	}

	beforeLength := store.Length()

	store.Delete(token)

	afterLength := store.Length()

	if afterLength != beforeLength-1 {
		t.Errorf("Failed to delete session \n Before: %v \n After: %v", beforeLength, afterLength-1)
	}
}

func TestSessionExists(t *testing.T) {
	store := NewDefaultSessionManager()

	token, err := store.NewSession()
	if err != nil {
		t.Errorf("Failed to create session")
	}
	exists := store.Exists(token)

	if exists != true {
		t.Error("context does not exist when it should")
	}

	b := make([]byte, 64)
	_, err = rand.Read(b)
	if err != nil {
		t.Errorf("Could not generate random string")
	}

	stringToken := base64.URLEncoding.EncodeToString(b)
	badToken := SessionToken(stringToken)

	exists = store.Exists(badToken)

	if exists == true {
		t.Error("context exists when it should not")
	}
}

func TestSessionExpired(t *testing.T) {
	store := NewDefaultSessionManager()
	token, err := store.NewSession()
	if err != nil {
		t.Errorf("Failed to create session")
	}

	ctx, _ := store.Get(token)
	expireTime := time.Now().Add(time.Second)
	ctx.expireTime = expireTime
	ctx.Client = &Client{
		mutex: &sync.Mutex{},
		open:  false,
	}

	time.Sleep(time.Second)

	if !ctx.SessionExpired() {
		t.Errorf("Session did not expire after period of disconnection")
	}
}

func TestSetClient(t *testing.T) {
	store := NewDefaultSessionManager()
	token, err := store.NewSession()
	if err != nil {
		t.Errorf("Failed to create session")
	}

	err = store.SetClient(token, &Client{})

	if err != nil {
		t.Error(err)
	}
}

func TestGetByClient(t *testing.T) {
	store := NewDefaultSessionManager()
	token, err := store.NewSession()
	if err != nil {
		t.Errorf("Failed to create session")
	}

	client := &Client{}

	err = store.SetClient(token, client)
	if err != nil {
		t.Error(err)
	}

	session, err := store.GetByClient(client)
	if err != nil {
		t.Errorf("Failed to get client")
	}

	if session.Client != client {
		t.Errorf("Client does not match client! %v != %v", client, session.Client)
	}

}

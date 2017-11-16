package carrot

import (
	"crypto/rand"
	"encoding/base64"
	"testing"
	"time"
)

func TestRefreshExpiryTime(t *testing.T) {
	refreshExpiryTime()
}

func TestDefaultSessionManagerNewSessionGet(t *testing.T) {
	store := NewDefaultSessionManager()

	token, _, err := store.NewSession()
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

	token, _, err := store.NewSession()
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

	token, _, err := store.NewSession()
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

	token, _, err := store.NewSession()
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
	token, _, err := store.NewSession()
	if err != nil {
		t.Errorf("Failed to create session")
	}

	ctx, _ := store.Get(token)
	expireTime := time.Now().Add(time.Second)
	ctx.expireTime = expireTime

	time.Sleep(time.Second)

	if !ctx.sessionDurationExpired() {
		t.Errorf("Session did not expire after period of disconnection")
	}
}

func TestPrimaryDeviceAssignmentAndRetrieval(t *testing.T) {
	store := NewDefaultSessionManager()
	_, session, err := store.NewSession()
	if err != nil {
		t.Error(err)
	}
	if session.isPrimaryDevice() {
		t.Errorf("The session should not be marked as a primary device")
	}
	_, err = store.GetPrimaryDeviceToken()
	if err == nil {
		t.Errorf("No primary device should have been retrieved because one doesn't exist yet")
	}
	session.primaryDevice = true
	if !session.isPrimaryDevice() {
		t.Errorf("The session should have been marked as a primary device")
	}
	_, err = store.GetPrimaryDeviceToken()
	if err != nil {
		t.Error(err)
	}
}

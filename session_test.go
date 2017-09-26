package buddy

import (
	"crypto/rand"
	"encoding/base64"
	"testing"
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
		t.Error("Context was not received")
	}
}

func TestSessionDelete(t *testing.T) {
	store := NewDefaultSessionManager()

	token, err := store.NewSession()
	if err != nil {
		t.Errorf("Failed to create session")
	}
	store.Delete(token)
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

func TestSessionLength(t *testing.T) {
	store := NewDefaultSessionManager()
	token, err := store.NewSession()
	if err != nil {
		t.Errorf("Failed to create session")
	}

	expectedLength := 1
	actualLength := store.Length()

	if expectedLength != actualLength {
		t.Errorf("Length should have been %v but was %v", expectedLength, actualLength)
	}

	store.Delete(token)

	expectedLength = 0
	actualLength = store.Length()

	if expectedLength != actualLength {
		t.Errorf("Length should have been %v but was %v", expectedLength, actualLength)
	}

}

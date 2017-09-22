package buddy

import (
	"crypto/rand"
	"encoding/base64"
	"testing"
)

var (
	store = NewDefaultSessionManager()
)

func TestDefaultSessionManagerNewSessionGet(t *testing.T) {
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
	token, err := store.NewSession()
	if err != nil {
		t.Errorf("Failed to create session")
	}

	ctx, _ := store.Get(token)

	ctx["testVal"] = true

	modifiedCtx, _ := store.Get(token)

	if modifiedCtx["testVal"] != true {
		t.Error("context did not save after setting value")
	}
}

func TestSessionDelete(t *testing.T) {
	token, err := store.NewSession()
	if err != nil {
		t.Errorf("Failed to create session")
	}
	store.Delete(token)
}

func TestSessionExists(t *testing.T) {
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

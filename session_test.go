package buddy

import (
	"testing"
)

var (
	store = NewDefaultSessionManager()
)

func TestDefaultSessionManagerNewSessionGet(t *testing.T) {
	token := store.NewSession()

	_, err := store.Get(token)
	if err != nil {
		t.Errorf("Failed to get session %v", err)
	}
}

func TestContextPersistence(t *testing.T) {
	token := store.NewSession()

	ctx, _ := store.Get(token)

	ctx["testVal"] = true

	modifiedCtx, _ := store.Get(token)

	if modifiedCtx["testVal"] != true {
		t.Error("context did not save after setting value")
	}
}

func TestSessionDelete(t *testing.T) {
	token := store.NewSession()
	store.Delete(token)
}

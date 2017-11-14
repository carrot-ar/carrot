package carrot

import "testing"

func TestInSlice(t *testing.T) {
	items := []string{"potato", "radish", "carrot"}
	existsStr := "potato"
	doesNotExistStr := "beet"

	if !InSlice(existsStr, items) {
		t.Errorf("%v was in items when it shouldn't have been", existsStr)
	}

	if InSlice(doesNotExistStr, items) {
		t.Errorf("%v was in items when it shouldn't have been", doesNotExistStr)
	}
}

func TestGenerateUUID(t *testing.T) {
	_, err := generateUUID()
	if err != nil {
		t.Error(err)
	}
}

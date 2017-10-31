package carrot

import (
	"testing"
)

func TestRun(t *testing.T) {
	err := Run()
	if err != nil {
		t.Fatal("Failed to go through bootup sequence")
	}
}

package carrot

import (
	"testing"
)

func TestControllerFactory(t *testing.T) {
	_, err := NewController(TestController{}, false)
	if err != nil {
		t.Error(err)
	}
	_, err = NewController(TestStreamController{}, true)
	if err != nil {
		t.Error(err)
	}
}

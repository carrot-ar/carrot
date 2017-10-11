package buddy

import (
	"testing"
)


type SphereController struct{}

func TestControllerRouteLookup(t *testing.T) {
	Add("place_sphere", SphereController{}, "Place")
	actual := Lookup("place_sphere")

	expected := Route{
		controller: SphereController{},
		function: "Place",
		persist: false,
	}

	if actual != expected {
		t.Errorf("Routes do not match: %v != %v", actual, expected)
	}
}

/*
type DrawingController struct{}

func TestStreamControllerRouteLookup(t *testing.T) {
	Add("draw", DrawingController{}, "Draw")
	actual := Lookup("draw")

	expected := Route{
		controller: DrawingController{},
		function: "Draw",
		persist: true,
	}

	if actual != expected {
		t.Errorf("Routes do not match: %v != %v", actual, expected)
	}
}
*/
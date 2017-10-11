package carrot

import (
	"testing"
)

func TestControllerRouteLookup(t *testing.T) {
	actual := Lookup("place_sphere")

	expected := Route{
		controller: "Sphere",
		function: "Place",
		persist: false,
	}

	if actual != expected {
		t.Errorf("Routes do not match: %v != %v", actual, expected)
	}
}

func TestStreamControllerRouteLookup(t *testing.T) {
	actual := Lookup("draw")

	expected := Route{
		controller: "Drawing",
		function: "Draw",
		persist: true,
	}

	if actual != expected {
		t.Errorf("Routes do not match: %v != %v", actual, expected)
	}
}
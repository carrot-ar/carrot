package carrot

import (
	"testing"
)

func TestDispatchRequest(t *testing.T) {
	d := NewDispatcher()

	token1, route1, req1, err := getTokenRouteAndRequestForTest(endpoint1)
	if err != nil {
		t.Error(err)
	}

	key1 := getCacheKey(token1, route1.controller)
	if d.cachedControllers.Exists(key1) {
		t.Fatalf("The DPC is already in the cached controllers list")
	}

	if d.cachedControllers.Length() != 0 {
		t.Fatalf("No (zero) cached controllers should be stored yet")
	}

	err = d.dispatchRequest(route1, req1)
	if err != nil {
		t.Fatal(err)
	}

	if !d.cachedControllers.Exists(key1) {
		t.Fatalf("The DPC was unsuccessfully stored in the cached controllers list")
	}
	if d.cachedControllers.Length() != 1 {
		t.Fatalf("Only 1 (one) cached controller should be stored")
	}

	err = d.dispatchRequest(route1, req1)
	if err != nil {
		t.Fatal(err)
	}

	if d.cachedControllers.Length() != 1 { //should return int of 1 (using existing entry in list)
		t.Fatalf("An existing cached controller should have been used")
	}

	_, route2, req2, err := getTokenRouteAndRequestForTest(endpoint1)
	if err != nil {
		t.Error(err)
	}

	err = d.dispatchRequest(route2, req2)
	if err != nil {
		t.Fatal(err)
	}

	if d.cachedControllers.Length() != 2 { //should return int of 2
		t.Fatalf("Only 2 (two) cached controllers should be stored")
	}
}

package carrot

import (
	"fmt"
	"testing"
	"time"
)

type TestDispatcherController struct{}

func (tdc *TestDispatcherController) Print(req *Request, broadcast *Broadcast) {
	fmt.Println("The dispatchingRequest method worked because this controller is speaking :D!")
}

func getTokenRouteAndRequestForTest(endpoint string) (SessionToken, *Route, *Request, error) {
	ss := NewDefaultSessionManager()
	token, err := ss.NewSession()
	if err != nil {
		return "", nil, nil, err
	}
	sesh, err := ss.Get(token)
	if err != nil {
		return "", nil, nil, err
	}
	req := &Request{
		SessionToken: sesh.Token,
		metrics:      make([]time.Time, MetricCount),
		endpoint:     endpoint,
	}
	route, err := Lookup(req.endpoint)
	if err != nil {
		return "", nil, nil, err
	}
	return token, route, req, nil
}

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

	d.dispatchRequest(route1, req1)

	if !d.cachedControllers.Exists(key1) {
		t.Fatalf("The DPC was unsuccessfully stored in the cached controllers list")
	}
	if d.cachedControllers.Length() != 1 {
		t.Fatalf("Only 1 (one) cached controller should be stored")
	}

	d.dispatchRequest(route1, req1)

	if d.cachedControllers.Length() != 1 { //should return int of 1 (using existing entry in list)
		t.Fatalf("An existing cached controller should have been used")
	}

	_, route2, req2, err := getTokenRouteAndRequestForTest(endpoint1)
	if err != nil {
		t.Error(err)
	}

	d.dispatchRequest(route2, req2)

	if d.cachedControllers.Length() != 2 { //should return int of 2
		t.Fatalf("Only 2 (two) cached controllers should be stored")
	}
}

package carrot

import (
	"testing"
)

func TestMethodInvocation(t *testing.T) {	
	_, _, req1, err := getTokenRouteAndRequestForTest(endpoint1)
	if err != nil {
		t.Error(err)
	}

	_, _, req2, err := getTokenRouteAndRequestForTest(endpoint2)
	if err != nil {
		t.Error(err)
	}

	_, _, req3, err := getTokenRouteAndRequestForTest(endpoint2)
	if err != nil {
		t.Error(err)
	}

	_, _, req4, err := getTokenRouteAndRequestForTest(endpoint2)
	if err != nil {
		t.Error(err)
	}

	d := NewDispatcher()
	go d.Run()
	d.requests <- req1
	d.requests <- req2
	d.requests <- req3
	d.requests <- req4
}

func TestInvalidMethodInvocation(t *testing.T) {
	_, _, req, err := getTokenRouteAndRequestForTest(endpoint3)
	if err != nil {
		t.Error(err)
	}

	rt, _ := Lookup(endpoint3)

	c, err := NewController(rt.Controller(), false)
	if err != nil {
		t.Errorf("Failed to make controller")
	}
	
	err = c.Invoke(rt, req)
	if err == nil {
		t.Errorf("Method invocation did not capture invalid method and probably crashed.")
	}
}

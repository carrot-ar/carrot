package carrot

import (
	"testing"
	"time"
	"fmt"
)

type TestDispatcherController struct {}

func (tdc *TestDispatcherController) Print(req *Request, broadcast *Broadcast) {
	fmt.Println("The dispatchingRequest method worked because this controller is speaking!")
}

func TestDispatchRequest(t *testing.T) {
	Add("test1", TestDispatcherController{}, "Print", false)

	ss := NewDefaultSessionManager()

	token1, err1 := ss.NewSession()
	if err1 != nil {
		fmt.Println(err1)
	}
	sesh1, err2 := ss.Get(token1)
	if err2 != nil {
		fmt.Println(err2)
	}

	req1 := &Request{
		sessionToken: sesh1.Token,
		metrics:      make([]time.Time, MetricCount),
		endpoint:     "test1",
	}

	token2, err3 := ss.NewSession()
	if err3 != nil {
		fmt.Println(err1)
	}
	sesh2, err4 := ss.Get(token2)
	if err4 != nil {
		fmt.Println(err2)
	}

	req2 := &Request{
		sessionToken: sesh2.Token,
		metrics:      make([]time.Time, MetricCount),
		endpoint:     "test2",
	}

	d := NewDispatcher()

	route1, err := Lookup(req1.endpoint)
	if err != nil {
		fmt.Println(err)
	}

	route2, err := Lookup(req2.endpoint)
	if err != nil {
		fmt.Println(err)
	}

	if (!d.cachedControllers.Exists(sesh1.Token)) {
		fmt.Println("The DPC is not yet in the cached controllers list")
	}
	fmt.Println(d.cachedControllers.Length()) //should return int of 0

	d.dispatchRequest(route1, req1)

	if (d.cachedControllers.Exists(sesh1.Token)) {
		fmt.Println("The DPC was successfully stored in the cached controllers list")
	}
	fmt.Println(d.cachedControllers.Length()) //should return int of 1

	d.dispatchRequest(route1, req1)

	fmt.Println(d.cachedControllers.Length()) //should return int of 1 (using existing entry in list)

	d.dispatchRequest(route2, req2)	
	
	fmt.Println(d.cachedControllers.Length()) //should return int of 2
}

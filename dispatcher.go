package carrot

import (
	"fmt"
)

const (
	doCacheControllers bool = true
)

type Dispatcher struct {
	cachedControllers *CachedControllersList
	requests    chan *Request
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		cachedControllers: NewCachedControllersList(),
		requests:    make(chan *Request, 256),
	}
}

func (dp *Dispatcher) dispatchRequest(route *Route, req *Request) {
	req.AddMetric(DispatchRequestStart)
	if doCacheControllers {	//used to be "if route.persist"
		token := req.sessionToken
		if exists := dp.cachedControllers.Exists(token); !exists {
			c, err := NewController(route.Controller(), doCacheControllers) //send to controller factory with stream identifier
			if err != nil {
				fmt.Println(err)
			}
			dp.cachedControllers.Add(token, c)
		}
		sc := dp.cachedControllers.Get(token)

		err := sc.Invoke(route, req) //send request to controller
		if err != nil {
			req.err = err
		}

	} else { //route leads to event controller
		c, err := NewController(route.Controller(), doCacheControllers) //send to controller factory with event identifier
		if err != nil {
			fmt.Println(err)
		}
		err = c.Invoke(route, req) //send request to controller
		req.err = err
	}
}

func (dp *Dispatcher) Run() {
	for {
		select {
		case req := <-dp.requests:
			req.AddMetric(DispatchLookupStart)
			route, err := Lookup(req.endpoint)
			if err != nil {
				fmt.Println(err)
			}
			req.AddMetric(DispatchLookupEnd)

			req.AddMetric(DispatchRequestStart)
			dp.dispatchRequest(route, req)
			req.AddMetric(DispatchRequestEnd)
		default:
			// fmt.Println("dispatcher Run() did something bad")
		}
	}
}

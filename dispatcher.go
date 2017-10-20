package carrot

import (
	"fmt"
)

type Dispatcher struct {
	openStreams *OpenStreamsList
	requests    chan *Request
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		openStreams: NewOpenStreamsList(),
		requests:    make(chan *Request, 256),
	}
}

func (dp *Dispatcher) dispatchRequest(route *Route, req *Request) {
	req.AddMetric(DispatchRequestStart)
	if route.persist {
		token := req.sessionToken
		if exists := dp.openStreams.Exists(token); !exists {
			c1, err := NewController(route.Controller(), true) //send to controller factory with stream identifier
			if err != nil {
				fmt.Println(err)
			}
			dp.openStreams.Add(token, c1)
		}
		sc := dp.openStreams.Get(token)

		err := sc.Invoke(route, req) //send request to controller
		if err != nil {
			req.err = err
		}

	} else { //route leads to event controller
		c, err := NewController(route.Controller(), false) //send to controller factory with event identifier
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
			route := Lookup(req.endpoint)
			req.AddMetric(DispatchLookupEnd)

			req.AddMetric(DispatchRequestStart)
			dp.dispatchRequest(&route, req)
			req.AddMetric(DispatchRequestEnd)
		default:
			// fmt.Println("dispatcher Run() did something bad")
		}
	}
}

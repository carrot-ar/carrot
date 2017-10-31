package carrot

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	doCacheControllers bool = true
	maxNumCachedControllers = 256
)

type Dispatcher struct {
	cachedControllers *CachedControllersList
	requests          chan *Request
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		cachedControllers: NewCachedControllersList(),
		requests:          make(chan *Request, 256),
	}
}

func (dp *Dispatcher) dispatchRequest(route *Route, req *Request) {
	req.AddMetric(DispatchRequestStart)
	if doCacheControllers { //used to be "if route.persist"
		token := req.sessionToken
		key := getCacheKey(token, route.controller)
		if exists := dp.cachedControllers.Exists(key); !exists {
			c, err := NewController(route.Controller(), doCacheControllers) //send to controller factory with stream identifier
			if err != nil {
				fmt.Println(err)
			}
			dp.cachedControllers.Add(key, c)
		}
		sc := dp.cachedControllers.Get(key)

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
			//delete controllers that haven't been used recently
			if (dp.cachedControllers.Length() > maxNumCachedControllers) {
				dp.cachedControllers.DeleteOldest()
				fmt.Printf("a controller has been deleted, num of controllers left: %v \n", dp.cachedControllers.lru.Len())
			}
		}
	}
}

func getCacheKey(token SessionToken, controller ControllerType) string {
	c1 := reflect.TypeOf(controller)
	c2 := strings.SplitAfter(c1.String(), ".")
	tmp := []string{string(token), c2[1]}
	key := strings.Join(tmp, ".")
	fmt.Println(key)
	return key
	// will look like "token.controllerType"
}

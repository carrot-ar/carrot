package carrot

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

const (
	doCacheControllers      bool = true
	maxNumCachedControllers      = 256
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
		token := req.SessionToken
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
			log.WithFields(log.Fields{
				"route" : fmt.Sprintf("%v#%v", reflect.TypeOf(route.controller).Name(), route.function),
				"session_token" : req.SessionToken,
			}).Debug("dispatching request")
		default:
			//delete controllers that haven't been used recently
			if dp.cachedControllers.Length() > maxNumCachedControllers {
				dp.cachedControllers.DeleteOldest()
				log.WithField("cache_size", dp.cachedControllers.lru.Len()).Debug("deleting least recently used controller")
			}
		}
	}
}

func getCacheKey(token SessionToken, controller ControllerType) string {
	c1 := reflect.TypeOf(controller)
	c2 := strings.SplitAfter(c1.String(), ".")
	tmp := []string{string(token), c2[1]}
	key := strings.Join(tmp, ".")
	return key
	// will look like "token.controllerType"
}

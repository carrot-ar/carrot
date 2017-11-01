package carrot

import (
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

func (dp *Dispatcher) dispatchRequest(route *Route, req *Request) error {
	req.AddMetric(DispatchRequestStart)
	if doCacheControllers { //used to be "if route.persist"
		token := req.SessionToken
		key := getCacheKey(token, route.controller)
		if exists := dp.cachedControllers.Exists(key); !exists {
			c, err := NewController(route.Controller(), doCacheControllers) //send to controller factory with stream identifier
			if err != nil {
				return err
			}
			dp.cachedControllers.Add(key, c)
		}
		sc, err := dp.cachedControllers.Get(key)
		reqErr := sc.Invoke(route, req) //send request to controller
		if reqErr != nil {
			req.err = reqErr
		}
		return err
	} else { //route leads to event controller
		c, err := NewController(route.Controller(), doCacheControllers) //send to controller factory with event identifier
		if err != nil {
			return err
		}
		err = c.Invoke(route, req) //send request to controller
		req.err = err
	}
	return nil
}

func (dp *Dispatcher) Run() {
	for {
		select {
		case req := <-dp.requests:
			req.AddMetric(DispatchLookupStart)
			route, err := Lookup(req.endpoint)
			if err != nil {
				log.Error(err)
			}
			req.AddMetric(DispatchLookupEnd)

			req.AddMetric(DispatchRequestStart)
			err = dp.dispatchRequest(route, req)
			if err != nil {
				log.Error(err)
			}
			req.AddMetric(DispatchRequestEnd)
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

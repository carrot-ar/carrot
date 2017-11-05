package carrot

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"reflect"
	"strings"
)

const (
	doCacheControllers               bool = true
	maxNumCachedControllers               = 4096
	maxNumDispatcherIncomingRequests      = 4096
)

type Dispatcher struct {
	cachedControllers *CachedControllersList
	requests          chan *Request
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		cachedControllers: NewCachedControllersList(),
		requests:          make(chan *Request, maxNumDispatcherIncomingRequests),
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
			if len(dp.requests) > int(math.Floor(maxNumDispatcherIncomingRequests*0.90)) {
				log.WithFields(log.Fields{
					"size":   len(dp.requests),
					"module": "dispatcher"}).Warn("input channel is at or above 90% capacity!")
			}

			if len(dp.requests) == maxNumDispatcherIncomingRequests {
				log.WithFields(log.Fields{
					"size":   len(dp.requests),
					"module": "dispatcher"}).Error("input channel is full!")
			}

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

			log.WithFields(log.Fields{
				"route":         fmt.Sprintf("%v#%v", reflect.TypeOf(route.controller).Name(), route.function),
				"session_token": req.SessionToken,
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

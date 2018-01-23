package carrot

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"reflect"
	"strings"
)

const (
	doCacheControllers               bool = true // determines whether controllers should be cached for streaming purposes
	maxNumCachedControllers               = 4096
	maxNumDispatcherIncomingRequests      = 16384
)

type Dispatcher struct {
	cachedControllers *CachedControllersList
	requests          chan *Request
	logger            *log.Entry
}

// NewDispatcher initializes a new instance of the Dispatcher struct.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		cachedControllers: NewCachedControllersList(),
		requests:          make(chan *Request, maxNumDispatcherIncomingRequests),
		logger:            log.WithField("module", "dispatcher"),
	}
}

// dispatchRequest sends requests to their specified controllers.
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

// Run establishes request routes, dispatches requests, and logs information.
func (dp *Dispatcher) Run() {
	for {
		select {
		case req := <-dp.requests:
			if len(dp.requests) > int(math.Floor(maxNumDispatcherIncomingRequests*0.90)) {
				dp.logger.WithField("buf_size", len(dp.requests)).Warn("input channel is at or above 90% capacity!")
			}

			if len(dp.requests) == maxNumDispatcherIncomingRequests {
				dp.logger.WithField("buf_size", len(dp.requests)).Error("input channel is full!")
			}

			req.AddMetric(DispatchLookupStart)
			route, err := Lookup(req.endpoint)
			if err != nil {
				dp.logger.Error(err)
				break
			}
			req.AddMetric(DispatchLookupEnd)

			req.AddMetric(DispatchRequestStart)
			err = dp.dispatchRequest(route, req)
			if err != nil {
				dp.logger.Error(err)
			}
			req.AddMetric(DispatchRequestEnd)

			dp.logger.WithFields(log.Fields{
				"route":         fmt.Sprintf("%v#%v", reflect.TypeOf(route.controller).Name(), route.function),
				"session_token": req.SessionToken,
			}).Debug("dispatching request")
		default:
			//delete controllers that haven't been used recently
			if dp.cachedControllers.Length() > maxNumCachedControllers {
				dp.cachedControllers.DeleteOldest()
				dp.logger.WithField("cache_size", dp.cachedControllers.lru.Len()).Debug("deleting least recently used controller")
			}
		}
	}
}

// getCacheKey returns a key for a controller type since only one instance of each controller type is cached.
func getCacheKey(token SessionToken, controller ControllerType) string {
	c1 := reflect.TypeOf(controller)
	c2 := strings.SplitAfter(c1.String(), ".")
	tmp := []string{string(token), c2[1]}
	key := strings.Join(tmp, ".")
	return key
	// will look like "token.controllerType"
}

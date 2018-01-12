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
	requests          chan *CContext
	logger            *log.Entry
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		cachedControllers: NewCachedControllersList(),
		requests:          make(chan *CContext, maxNumDispatcherIncomingRequests),
		logger:            log.WithField("module", "dispatcher"),
	}
}

func (dp *Dispatcher) dispatchRequest(route *Route, ctx *CContext) error {
	if doCacheControllers { //used to be "if route.persist"
		token := ctx.Session().Token
		key := getCacheKey(token, route.controller)
		if exists := dp.cachedControllers.Exists(key); !exists {
			c, err := NewController(route.Controller(), doCacheControllers) //send to controller factory with stream identifier
			if err != nil {
				return err
			}
			dp.cachedControllers.Add(key, c)
		}
		sc, err := dp.cachedControllers.Get(key)
		reqErr := sc.Invoke(route, ctx) //send request to controller
		if reqErr != nil {
			ctx.error = reqErr
		}
		return err
	} else { //route leads to event controller
		c, err := NewController(route.Controller(), doCacheControllers) //send to controller factory with event identifier
		if err != nil {
			return err
		}
		err = c.Invoke(route, ctx) //send request to controller
		ctx.error = err
	}
	return nil
}

func (dp *Dispatcher) Run() {
	for {
		select {
		case ctx := <-dp.requests:
			if len(dp.requests) > int(math.Floor(maxNumDispatcherIncomingRequests*0.90)) {
				dp.logger.WithField("buf_size", len(dp.requests)).Warn("input channel is at or above 90% capacity!")
			}

			if len(dp.requests) == maxNumDispatcherIncomingRequests {
				dp.logger.WithField("buf_size", len(dp.requests)).Error("input channel is full!")
			}

			route, err := Lookup(ctx.Request().Endpoint)
			if err != nil {
				dp.logger.Error(err)
				break
			}
			//req.AddMetric(DispatchLookupEnd)

			//req.AddMetric(DispatchRequestStart)
			err = dp.dispatchRequest(route, ctx)
			if err != nil {
				dp.logger.Error(err)
			}
			//req.AddMetric(DispatchRequestEnd)

			dp.logger.WithFields(log.Fields{
				"route":         fmt.Sprintf("%v#%v", reflect.TypeOf(route.controller).Name(), route.function),
				"session_token": ctx.Session().Token,
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

func getCacheKey(token SessionToken, controller ControllerType) string {
	c1 := reflect.TypeOf(controller)
	c2 := strings.SplitAfter(c1.String(), ".")
	tmp := []string{string(token), c2[1]}
	key := strings.Join(tmp, ".")
	return key
	// will look like "token.controllerType"
}

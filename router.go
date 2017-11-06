package carrot

import (
	"fmt"
	"sync"
)

var (
	routerConfig         = config.Router
	routeDelimiter       = routerConfig.RouteDelimiter
	streamIdentifier     = routerConfig.StreamIdentifier
	controllerIdentifier = routerConfig.ControllerIdentifier
)

type RoutingTable map[string]Route

var (
	oneRouter      sync.Once
	routerInstance Router
)

type Router interface {
	addRoute(Endpoint, *Route)
	get(Endpoint) (*Route, error)
	Range(func(key, value interface{}) bool)
	Length() int
}

type DefaultRouter struct {
	routingTable *sync.Map
	length       int
	mutex        *sync.Mutex
}

func getRouter() Router {
	oneRouter.Do(func() {
		routerInstance = &DefaultRouter{
			routingTable: &sync.Map{},
			length:       0,
			mutex:        &sync.Mutex{},
		}
	})

	return routerInstance
}

func (r *DefaultRouter) addRoute(endpoint Endpoint, route *Route) {
	r.routingTable.Store(endpoint, route)
	r.mutex.Lock()
	r.length += 1
	r.mutex.Unlock()
}

func (r *DefaultRouter) get(endpoint Endpoint) (*Route, error) {
	route, ok := r.routingTable.Load(endpoint)
	if !ok {
		return nil, fmt.Errorf("route does not exist")
	}

	return route.(*Route), nil
}

func (r *DefaultRouter) Range(f func(key, value interface{}) bool) {
	r.routingTable.Range(f)
}

func (r *DefaultRouter) Length() int {
	r.mutex.Lock()
	length := r.length
	r.mutex.Unlock()

	return length
}

func Lookup(endpoint string) (*Route, error) {
	router := getRouter()
	route, err := router.get(Endpoint(endpoint))
	if err != nil {
		return nil, err
	}

	return route, nil
}

func Add(endpoint string, controller ControllerType, function string, persist bool) {
	getRouter().addRoute(Endpoint(endpoint),
		&Route{controller: controller,
			function: function,
			persist:  persist,
		})
}

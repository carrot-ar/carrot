package buddy

import (
	"reflect"
)

/*
	HOW THIS WORKS!
	All controllers will contain a field referencing the "Controller" type. Then,
	to implement controllers, the developer will implement "Actions" which receive
    a request and a (in the future) a responder type. To invoke these actions,
	the 'Invoke' method will be used which receives a route and a request.
	The Invoke function will perform reflection to find the proper function within
	the controller, then call it.
	A responder object will be used for the developer to decide which group of users
	receive the broadcast (to be implemented with broadcast group middleware
 */


/*
	AppController. This is a default controller on which simple routes can be made.
	Adding a route such as "place_ball => Controller_PlaceBall" will route to the
	PlaceBall function within the controller
 */
type AppController struct {
	persist bool
	Controller interface{}
	//reqBuffer chan *buddy.Request
	// add a responder here for responding to all, groups, individual etc"
}

func (c *AppController) Persist(p bool) {
	c.persist = p
}

/*
		Reflect on the controller and find the correct function to call, then call it
*/
func (c *AppController) Invoke(route Route, req *Request) {

	req.AddMetric(ControllerInvocation)

	// Create a new Value pointer representing the controller type
	ptr := reflect.New(reflect.TypeOf(c.Controller))

	// Look at that value then call the correct method
	method := ptr.MethodByName(route.Function())

	if method.IsValid() {
		args := []reflect.Value{reflect.ValueOf(req)}
		method.Call(args)
	}
}
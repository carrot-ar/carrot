package carrot

import (
	"fmt"
	"reflect"
)

type AppController struct {
	persist    bool
	Controller ControllerType
	//reqBuffer chan *buddy.Request
	// add a responder here for responding to all, groups, individual etc"
	broadcast *Broadcast
}

type ControllerType interface{}

func (c *AppController) Persist(p bool) {
	c.persist = p
}

// Invoke reflects on the controller to find the correct function to call and then calls it.
func (c *AppController) Invoke(route *Route, req *Request) error {
	req.AddMetric(MethodReflectionStart)

	ptr := c.Controller.(reflect.Value).Type()

	method, ok := ptr.MethodByName(route.Function())

	if ok != true {
		return fmt.Errorf("method lookup failed")
	}

	if method.Func.IsValid() {
		args := []reflect.Value{c.Controller.(reflect.Value), reflect.ValueOf(req), reflect.ValueOf(c.broadcast)}

		req.AddMetric(MethodReflectionEnd)
		req.AddMetric(ControllerMethodStart)

		method.Func.Call(args)
		req.AddMetric(ControllerMethodEnd)
		//req.End()
	} else {
		return fmt.Errorf("invalid method called")
	}

	return nil
}

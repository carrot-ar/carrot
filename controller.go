package buddy

import (
	"fmt"
	"log"
	"reflect"
)

type AppController struct {
	persist    bool
	Controller ControllerType
	//reqBuffer chan *buddy.Request
	// add a responder here for responding to all, groups, individual etc"
	responder *Responder
}

type ControllerType interface{}

func (c *AppController) Persist(p bool) {
	c.persist = p
}

/*
	Reflect on the controller and find the correct function to call, then call it
*/
func (c *AppController) Invoke(route *Route, req *Request) {

	req.AddMetric(ControllerInvocation)

	ptr := c.Controller.(reflect.Value).Type()

	method, ok := ptr.MethodByName(route.Function())

	if ok != true {
		fmt.Println("the method is not ok!")
	}

	if method.Func.IsValid() {
		args := []reflect.Value{c.Controller.(reflect.Value), reflect.ValueOf(req), reflect.ValueOf(c.responder)}
		method.Func.Call(args)
	} else {
		log.Printf("error: invalid method called")
	}
}

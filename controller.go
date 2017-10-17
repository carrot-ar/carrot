package buddy

import (
	"reflect"
	"log"
	"fmt"
)

type AppController struct {
	persist bool
	Controller ControllerType
	//reqBuffer chan *buddy.Request
	// add a responder here for responding to all, groups, individual etc"
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
	fmt.Println(route)

	// Create a new Value pointer representing the controller type
	ptr := reflect.New(reflect.TypeOf(c.Controller))
	fmt.Println(reflect.TypeOf(c.Controller))
	fmt.Println(ptr)
	fmt.Println("VALUE==========")
	val := reflect.ValueOf(&c.Controller).Elem()
	fmt.Println(val.Elem().Interface())
	fmt.Println(val.Elem().NumMethod())
	fmt.Println(val.Elem().MethodByName(route.Function()).Interface())

	// Look at that value then call the correct method
	method := ptr.MethodByName(route.Function())

	if method.IsValid() {
		args := []reflect.Value{reflect.ValueOf(req)}
		method.Call(args)
	} else {
		log.Printf("error: invalid method called")
	}
}
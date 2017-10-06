package controller

import(
	"github.com/senior-buddy/buddy"
	"github.com/senior-buddy/buddy/routes"
	"reflect"
	"fmt"
	"time"
)

const (
	/*
		Add controllers below this line
	 */
	APP_CONTROLLER = iota
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
type BaseController struct {
	persist bool
	Controller interface{}
	//reqBuffer chan *buddy.Request
	// add a responder here for responding to all, groups, individual etc"
}

func (c *BaseController) Persist(p bool) {
	c.persist = p
}

func (c *BaseController) Invoke(route routes.Route, req *buddy.Request) {
	/*
		Reflect on the controller and find the correct function to call, then call it
	 */
	start := time.Now()

	v := reflect.ValueOf(c.Controller)
	fmt.Printf("Value of controller: %v\n", v)

	fmt.Printf("Kind is %v and Ptr is %v\n", v.Type().Kind(), reflect.Ptr)

	// Make a new pointer to the type
	ptr := reflect.New(reflect.TypeOf(c.Controller))

	// create a temp version of the pointer
	temp := ptr.Elem()

	// Make that pointer the passed value
	temp.Set(v)

	fmt.Println(temp)

	// now for method lookup!
	method := ptr.MethodByName(route.Function())

	fmt.Println(method)
	if method.IsValid() {
		args := []reflect.Value{reflect.ValueOf(req)}
		method.Call(args)

	}

	fmt.Println(time.Now().Sub(start))
	 /*
	 // Get the Value or reflection interface of the controller
	fmt.Println("PRINTING CONTROLLER TYPE!")
	fmt.Println(reflect.TypeOf(c.Controller))
	fmt.Println(c.Controller)

	v := reflect.ValueOf(c.Controller)
	fmt.Println(v)

	// Get the function to call from the Function string provided by the route
	m := v.MethodByName(route.Function())

	// build up args using their reflection
	args := []reflect.Value{reflect.ValueOf(req)}

	// Call the function
	fmt.Println(v)
	fmt.Println(m)
	fmt.Println(args)
	m.Call(args)
	*/
}

// Base controller
//type ControllerType int

/*

// returns a new controller of controllerType
func New(t ControllerType) (Controller, error) {
	switch t {
	case APP_CONTROLLER:
		return &AppController{ persist: false }, nil
	default:
		return &AppController{ persist: false }, nil
	}
}

*/
package buddy

import (
	"fmt"
)

type Dispatcher struct {
	openStreams *OpenStreamsList
	requests 	chan *Request
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		openStreams:	NewOpenStreamsList(),
		requests:		make(chan *Request, 256),
	}
}

func (dp *Dispatcher) dispatchRequest(route *Route, req *Request) {
	if route.persist {
		token := req.sessionToken
		if exists := dp.openStreams.Exists(token); !exists {
			c, err := NewController(route.Controller()) //send controller string to controller factory
			if err != nil {
				fmt.Println(err)
			}	
			dp.openStreams.Add(token, c)
		}
		sc := dp.openStreams.Get(token)
		sc.Invoke(route, req) //send request to controller		
	} else { //route leads to event controller
		fmt.Println("CONTROLLER INFO!")
		fmt.Println(route)
		fmt.Println(route.function)
		fmt.Println(route.controller)
		c, err := NewController(route.Controller()) //send controller string to controller factory
		fmt.Println(c)

		if err != nil {
			fmt.Println(err)
		}		
		c.Invoke(route, req) //send request to controller
		
	}
}

func (dp *Dispatcher) Run() {
	for {
		select {
		case req := <- dp.requests:
			fmt.Println("HEY I GOT A FUCKING REQUEST")
			fmt.Println(req.endpoint)
			route := Lookup(req.endpoint)
			fmt.Println(route)
			dp.dispatchRequest(&route, req)
		default:
			// fmt.Println("dispatcher Run() did something bad")
		}
	}
}

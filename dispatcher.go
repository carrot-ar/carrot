package buddy

import (
	"fmt"
)

type Dispatcher struct {
	openStreams *OpenStreamsList
	requests    chan *Request
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		openStreams: NewOpenStreamsList(),
		requests:    make(chan *Request, 256),
	}
}

func (dp *Dispatcher) dispatchRequest(route *Route, req *Request) {
	if route.persist {
		fmt.Printf("\nopenStreams is empty...%v\n", dp.openStreams.IsEmpty())
		token := req.session.Token
		if exists := dp.openStreams.Exists(token); !exists {
			fmt.Printf("\ntoken is:	%v\n", token)
			c1, err := NewController(route.Controller(), true) //send to controller factory with stream identifier
			if err != nil {
				fmt.Println(err)
			}
			dp.openStreams.Add(token, c1)
			fmt.Println("added to streams list")
			c2 := dp.openStreams.Get(token)
			if c1 == c2 {
				fmt.Println("the add and get work!!")
			}
		}
		sc := dp.openStreams.Get(token)
		sc1 := dp.openStreams.Get(token)
		if sc == sc1 {
			fmt.Println("Get returns the same object")
		}
		sc.Invoke(route, req) //send request to controller
		// sc.Invoke(route, req) //send request to controller

	} else { //route leads to event controller
		c, err := NewController(route.Controller(), false) //send to controller factory with event identifier
		if err != nil {
			fmt.Println(err)
		}
		c.Invoke(route, req) //send request to controller

	}
}

func (dp *Dispatcher) Run() {
	for {
		select {
		case req := <-dp.requests:
			route := Lookup(req.endpoint)
			dp.dispatchRequest(&route, req)
		default:
			// fmt.Println("dispatcher Run() did something bad")
		}
	}
}

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
		token := req.session.Token
		if exists := dp.openStreams.Exists(token); !exists {
			dp.openStreams.Add(token)
		}
		// sc := dp.openStreams.Get(token)
		//reflect on method
		//send request to controller
	} else { //route leads to event controller
		//send controller string to controller factory
		c, err := NewController(route.Controller())
		if err != nil {
			fmt.Println(err)
		}		
		//send request to controller
		c.Invoke(route, req)
	}
}

func (dp *Dispatcher) Run() {
	for {
		select {
		case req := <- dp.requests:
			fmt.Println(req)
			route := Lookup(req.endpoint)
			dp.dispatchRequest(&route, req)
		default:
			// fmt.Println("dispatcher Run() did something bad")
		}
	}
}

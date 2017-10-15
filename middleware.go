package buddy

import (
	"log"
)

const (
	InputChannelSize  = 256
	OutputChannelSize = 256
)

/*
	Middlewares
*/
// func parseRequest(req *Request) {
// 	loggerMw.Print("I am going to parse a request!")
// }

func logger(req *Request) {
	log.Println("middleware: new request")
}

type MiddlewarePipeline struct {
	In          chan *Request
	Out         chan *Request
	middlewares []func(*Request)
	dispatcher *Dispatcher
}

func (mw *MiddlewarePipeline) Run() {
	go mw.dispatcher.Run()

	func() {
		for {
			select {
			case req := <-mw.In:
				//req.AddMetric(MiddlewareInput)
				for _, f := range mw.middlewares {
					f(req)
				}
				mw.dispatcher.requests <- req
				//req.AddMetric(MiddlewareOutput)
				//mw.Out <- req
			}
		}
	}()
}

func NewMiddlewarePipeline() *MiddlewarePipeline {
	// List of middleware functions
	mw := []func(*Request){logger}

	return &MiddlewarePipeline{
		In:          make(chan *Request, InputChannelSize),
		Out:         make(chan *Request, OutputChannelSize),
		middlewares: mw,
		dispatcher: NewDispatcher(),
	}
}

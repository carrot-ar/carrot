package carrot

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
	log.Printf("middleware: new event: tbd | payload: %v\n", string(req.message[:]))
}

type MiddlewarePipeline struct {
	In          chan *Request
	Out         chan *Request
	middlewares []func(*Request)
}

func (mw *MiddlewarePipeline) Run() {
	func() {
		for {
			select {
			case req := <-mw.In:
				//req.AddMetric(MiddlewareInput)
				for _, f := range mw.middlewares {
					f(req)
				}
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
	}
}

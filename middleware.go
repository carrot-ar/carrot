package carrot

import (
	"log"
	"time"
)

const (
	InputChannelSize = 256
)

var count int = 0
var rate float64 = 0

/*
	Middlewares
*/
// func parseRequest(req *Request) {
// 	loggerMw.Print("I am going to parse a request!")
// }

func logger(req *Request) error {
	//log.Println("middleware: new request")
	return nil
}

func discardBadRequest(req *Request) error {
	if req.err != nil {
		log.Printf("bad request: %s, ignoring...\n", req.err.Error())
		return req.err
	}
	return nil
}

type MiddlewarePipeline struct {
	In          chan *Request
	middlewares []func(*Request) error
	dispatcher  *Dispatcher
}

func (mw *MiddlewarePipeline) Run() {
	go mw.dispatcher.Run()

	func() {
		for {
			select {
			case req := <-mw.In:
				req.AddMetric(MiddlewareInput)
				var err error
				for _, f := range mw.middlewares {
					err = f(req)
					if err != nil {
						req.End()
						break
					}
					count++
				}
				if err == nil {
					mw.dispatcher.requests <- req
				}
				req.AddMetric(MiddlewareOutputToDispatcher)
			}
		}
	}()
}

func NewMiddlewarePipeline() *MiddlewarePipeline {
	// List of middleware functions
	mw := []func(*Request) error{logger, discardBadRequest}

	seconds := 0
	go func() {
		for {
			time.Sleep(time.Second)
			seconds++
			rate = float64(count) / float64(seconds)
			log.Printf("%v requests per second\n", rate)
			log.Printf("%v request count\n", count)
		}
	}()
	return &MiddlewarePipeline{
		In:          make(chan *Request, InputChannelSize),
		middlewares: mw,
		dispatcher:  NewDispatcher(),
	}
}

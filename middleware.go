package buddy

import (
	"log"
	"os"
	"time"
)

const (
	InputChannelSize  = 256
	OutputChannelSize = 256
)

var (
	loggerMw = log.New(os.Stdout, "buddy: ", log.Ltime)
)

/*
	Middlewares
*/
// func parseRequest(req *Request) {
// 	loggerMw.Print("I am going to parse a request!")
// }

func logger(req *Request) {
	end := time.Now()
	loggerMw.Printf("middleware: new event: tbd | elapsed time: %v | payload: %v",
		end.Sub(req.startTime), string(req.message[:]))
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
				for _, f := range mw.middlewares {
					f(req)
				}
				mw.Out <- req
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

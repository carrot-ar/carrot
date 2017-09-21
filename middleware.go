package buddy

import (
	"log"
	"os"
)

const (
	InputChannelSize  = 256
	OutputChannelSize = 256
)

var (
	loggerMw = log.New(os.Stdout, "buddy: ", log.Lmicroseconds)
)

/*
	Middlewares
*/
func logger(ctx *Context) {
	loggerMw.Print("New Event: <event name> | Payload: <payload> | Other Data: <other_data>")
}

type MiddlewarePipeline struct {
	In          chan Context
	Out         chan Context
	middlewares []func(*Context)
}

func (mw *MiddlewarePipeline) Run() {
	func() {
		for {
			select {
			case ctx := <-mw.In:
				for _, f := range mw.middlewares {
					f(&ctx)
				}
				mw.Out <- ctx
			}
		}
	}()
}

func NewMiddlewarePipeline() *MiddlewarePipeline {
	// List of middleware functions
	mw := []func(*Context){logger}

	return &MiddlewarePipeline{
		In:          make(chan Context, InputChannelSize),
		Out:         make(chan Context, OutputChannelSize),
		middlewares: mw,
	}
}

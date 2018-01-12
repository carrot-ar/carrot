package carrot

import (
	log "github.com/sirupsen/logrus"
	"math"
)

const (
	InputChannelSize = 4096
)

var count = 0

/*
	Middlewares
*/
// func parseRequest(req *Request, logger *log.Entry) {
// 	logger.Info("I am going to parse a request!")
// }

func logger(ctx *CContext, logger *log.Entry) error {

	//logger.WithField("session_token", ctx.SessionToken).Debug("new request")
	return nil
}

func discardBadRequest(ctx *CContext, logger *log.Entry) error {
	if ctx.Error() != nil {
		//logger.WithField("session_token", ctx.SessionToken).Errorf("invalid request: %v", ctx.err.Error())
		return ctx.Error()
	}

	return nil
}

type MiddlewarePipeline struct {
	In          chan *CContext
	middlewares []func(*CContext, *log.Entry) error
	dispatcher  *Dispatcher
	logger      *log.Entry
}

func (mw *MiddlewarePipeline) Run() {
	go mw.dispatcher.Run()
	func() {
		for {
			select {
			case ctx := <-mw.In:
				if len(mw.In) > int(math.Floor(InputChannelSize*0.90)) {
					mw.logger.WithField("buf_size", len(mw.In)).Warn("input channel is at or above 90% capacity!")
				}
				if len(mw.In) == InputChannelSize {
					mw.logger.WithField("buf_size", len(mw.In)).Warn("input channel is full!")
				}

				var err error
				for _, f := range mw.middlewares {
					err = f(ctx, mw.logger)
					if err != nil {
						break
					}
					count++
				}

				if err == nil {
					mw.dispatcher.requests <- ctx
				}

			}
		}
	}()
}

func NewMiddlewarePipeline() *MiddlewarePipeline {
	// middleware function index
	mw := []func(*CContext, *log.Entry) error{discardBadRequest, logger}

	return &MiddlewarePipeline{
		In:          make(chan *CContext, InputChannelSize),
		middlewares: mw,
		dispatcher:  NewDispatcher(),
		logger:      log.WithField("module", "middleware"),
	}
}

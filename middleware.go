package carrot

import (
	log "github.com/sirupsen/logrus"
)

const (
	InputChannelSize = 256
)

/*
	Middlewares
*/
// func parseRequest(req *Request) {
// 	loggerMw.Print("I am going to parse a request!")
// }

func logger(req *Request) error {
	log.WithFields(log.Fields{
		"session_token": req.SessionToken,
		"module": "middleware"}).Debug("new request")
	return nil
}

func discardBadRequest(req *Request) error {
	if req.err != nil {
		log.WithFields(log.Fields{
			"session_token": req.SessionToken,
			"module": "middleware"}).Errorf("invalid request: %v", req.err.Error())
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
	mw := []func(*Request) error{discardBadRequest, logger}

	return &MiddlewarePipeline{
		In:          make(chan *Request, InputChannelSize),
		middlewares: mw,
		dispatcher:  NewDispatcher(),
	}
}

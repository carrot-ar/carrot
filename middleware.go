package carrot

import (
	log "github.com/sirupsen/logrus"
	"math"
)

/*
	Middlewares
*/
// func parseRequest(req *Request) {
// 	loggerMw.Print("I am going to parse a request!")
// }

var count int

func logger(req *Request) error {

	log.WithFields(log.Fields{
		"session_token": req.SessionToken,
		"module":        "middleware"}).Debug("new request")

	return nil
}

func discardBadRequest(req *Request) error {
	if req.err != nil {
		log.WithFields(log.Fields{
			"session_token": req.SessionToken,
			"module":        "middleware"}).Errorf("invalid request: %v", req.err.Error())
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
				if len(mw.In) > int(math.Floor(float64(config.Middleware.InputChannelSize)*0.90)) {
					log.WithFields(log.Fields{
						"size":   len(mw.In),
						"module": "middleware"}).Warn("input channel is at or above 90% capacity!")
				}
				if len(mw.In) == config.Middleware.InputChannelSize {
					log.WithFields(log.Fields{
						"size":   len(mw.In),
						"module": "middleware"}).Error("input channel is full!")
				}
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
	mw := []func(*Request) error{discardBadRequest, logger}

	/*
		seconds := 0
		go func() {
			for {
				time.Sleep(time.Second)
				seconds++
				rate = float64(count) / float64(seconds)

				log.WithFields(log.Fields{
					"rps":    rate,
					"module": "middleware",
				}).Info("middleware metrics")

			}
		}()
	*/
	return &MiddlewarePipeline{
		In:          make(chan *Request, config.Middleware.InputChannelSize),
		middlewares: mw,
		dispatcher:  NewDispatcher(),
	}
}

package carrot

import (
	log "github.com/sirupsen/logrus"
	"math"
	"github.com/DataDog/datadog-go/statsd"
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

func logger(req *Request, logger *log.Entry) error {

	logger.WithField("session_token", req.SessionToken).Debug("new request")
	return nil
}

func discardBadRequest(req *Request, logger *log.Entry) error {
	if req.err != nil {
		logger.WithField("session_token", req.SessionToken).Errorf("invalid request: %v", req.err.Error())
		return req.err
	}

	return nil
}

type MiddlewarePipeline struct {
	In          chan *Request
	middlewares []func(*Request, *log.Entry) error
	dispatcher  *Dispatcher
	logger      *log.Entry
	statsd 		*statsd.Client
}

func (mw *MiddlewarePipeline) Run() {
	go mw.dispatcher.Run()
	func() {
		for {
			select {
			case req := <-mw.In:
				if len(mw.In) > int(math.Floor(InputChannelSize*0.90)) {
					mw.logger.WithField("buf_size", len(mw.In)).Warn("input channel is at or above 90% capacity!")
				}
				if len(mw.In) == InputChannelSize {
					mw.logger.WithField("buf_size", len(mw.In)).Warn("input channel is full!")
				}

				mw.statsd.Incr("carrot.middleware.request_rate", nil, 1)

				var err error
				for _, f := range mw.middlewares {
					err = f(req, mw.logger)
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
	// middleware function index
	mw := []func(*Request, *log.Entry) error{discardBadRequest, logger}

	logger := log.WithField("module", "middleware")
	c, err := statsd.New("127.0.0.1:8125")
	if err != nil {
		logger.Error(err)
	}


	return &MiddlewarePipeline{
		In:          make(chan *Request, InputChannelSize),
		middlewares: mw,
		dispatcher:  NewDispatcher(),
		logger:      logger,
		statsd: 		 c,
	}
}

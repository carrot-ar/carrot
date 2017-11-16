package carrot

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	MetricCount = 12

	RequestCreation              = iota // 1
	MiddlewareInput                     // 2
	MiddlewareOutputToDispatcher        // 3
	DispatchLookupStart                 // 4
	DispatchLookupEnd                   // 5
	DispatchRequestStart                // 6
	DispatchRequestEnd                  // 7
	MethodReflectionStart               // 8
	MethodReflectionEnd                 // 9
	ControllerMethodStart               // 10
	ControllerMethodEnd                 // 11
	ResponderInvocation
	ResponderElapsed
)

// request structs

type Request struct {
	SessionToken SessionToken
	endpoint     string
	Params       map[string]interface{}
	Offset       *offset
	data         []byte
	metrics      []time.Time
	err          error
}

func NewRequest(session *Session, message []byte) *Request { //returns error,
	m := make([]time.Time, MetricCount)

	req := Request{
		SessionToken: session.Token,
		metrics:      m,
		data:         message,
	}

	var md messageData
	err := json.Unmarshal(message, &md)

	if err == nil {
		err = validSession(session.Token, SessionToken(md.SessionToken))
	}

	req.err = err
	req.endpoint = md.Endpoint
	req.Params = md.Payload.Params
	req.Offset = md.Payload.Offset

	req.AddMetric(RequestCreation)

	return &req
}

// Add the time that a request is created to the request metric tracker
func (r *Request) AddMetric(stage int) {
	r.metrics[stage] = time.Now()
}

func (r *Request) End() {
	//logBenchmarks(r.metrics)
}

func logBenchmarks(metrics []time.Time) {
	var prev time.Time
	for i, metric := range metrics {
		if i != 0 {
			if i == 1 {
				fmt.Printf("%v => %v\n", i, time.Now().Sub(metric))
			} else {
				fmt.Printf("%v => %v\n", i, metric.Sub(prev))

			}
		}
		prev = metric
	}
}

func validSession(serverToken SessionToken, clientToken SessionToken) error {
	log.Infof("client: %v, server %v", clientToken, serverToken)

	if serverToken != clientToken {
		return fmt.Errorf("client-server token mismatch")
	}

	return nil
}

package carrot

import (
	"encoding/json"
	"fmt"
	"time"
	"log"
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
	sessionToken SessionToken
	endpoint     string
	Params       map[string]string
	Origin       location
	Offset       offset
	data         []byte
	metrics      []time.Time
	err          error
}

type location struct {
	Longitude float64
	Latitude  float64
	Altitude float64
}

type offset struct {
	X float64
	Y float64
	Z float64
}

func NewRequest(session *Session, message []byte) *Request { //returns error,
	m := make([]time.Time, MetricCount)

	req := Request{
		sessionToken: session.Token,
		metrics:      m,
		data:         message,
	}

	log.Println(string(message))

	var d requestData
	err := json.Unmarshal(message, &d)

	if err == nil {
		err = validSession(session.Token, SessionToken(d.SessionToken))
	}

	req.err = err
	req.endpoint = d.Endpoint
	req.Params = d.Payload.Params
	req.Origin = location(d.Origin)
	req.Offset = offset(d.Payload.Offset)

	req.AddMetric(RequestCreation)

	return &req
}

// Add the time that a request is created to the request metric tracker
func (r *Request) AddMetric(stage int) {
	r.metrics[stage] = time.Now()
}

func (r *Request) End() {
	logBenchmarks(r.metrics)
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
	if serverToken != clientToken {
		return fmt.Errorf("token mismatch %v != %v", serverToken, clientToken)
	}

	return nil
}

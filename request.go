package buddy

import (
	"encoding/json"
	"fmt"
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

type Request struct {
	sessionToken SessionToken
	endpoint     string
	Params       map[string]string
	data         []byte
	metrics      []time.Time
	err          error
}

type requestData struct {
	SessionToken string            `json:"session_token"`
	Endpoint     string            `json:"endpoint"`
	Params       map[string]string `json:"params"`
}

func NewRequest(session *Session, message []byte) *Request { //returns error,
	m := make([]time.Time, MetricCount)

	req := Request{
		sessionToken: session.Token,
		metrics:      m,
		data:         message,
	}

	var d requestData
	err := json.Unmarshal(message, &d)

	err = validSession(session.Token, SessionToken(d.SessionToken))

	req.err = err
	req.endpoint = d.Endpoint
	req.Params = d.Params

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

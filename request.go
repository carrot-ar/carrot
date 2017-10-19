package buddy

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	MetricCount = 5

	RequestCreation = iota
	MiddlewareInput
	MiddlewareOutputToDispatcher
	ControllerInvocation
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

	var d requestData
	err := json.Unmarshal(message, &d)

	err = validSession(session.Token, SessionToken(d.SessionToken))

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
			fmt.Printf("%v => %v\n", i, metric.Sub(prev))
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

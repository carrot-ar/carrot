package buddy

import (
	//"fmt"
	//"log"
	"time"
	"fmt"
)

const (
	MetricCount = 2

	RequestCreation = 0
	//MiddlewareInput
	//MiddlewareOutput
	ControllerInvocation = 1
	ResponderInvocation
	ResponderElapsed
)

type Request struct {
	session   *Session
	message   []byte
	metrics []time.Time
	endpoint	string
}

// Add the time that a request is created to the request metric tracker
func (r *Request) AddMetric(stage int) {
	r.metrics[stage] = time.Now()
}

func NewRequest(session *Session, message []byte) *Request {
	m := make([]time.Time, MetricCount)
	m[RequestCreation] = time.Now()

	return &Request{
		session:   session,
		message:   message,
		metrics: 	m,
		endpoint:	"",
	}
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
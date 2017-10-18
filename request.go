package buddy

import (
	//"fmt"
	//"log"
	"encoding/json"
	"fmt"
	"time"
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
	sessionToken SessionToken
	endpoint     string
	Params       map[string]string
	data         []byte
	metrics      []time.Time
}

type requestData struct {
	Endpoint string            `json:"endpoint"`
	Params   map[string]string `json:"params"`
}

// Add the time that a request is created to the request metric tracker
func (r *Request) AddMetric(stage int) {
	r.metrics[stage] = time.Now()
}

func NewRequest(session *Session, message []byte) *Request { //returns error,
	m := make([]time.Time, MetricCount)
	m[RequestCreation] = time.Now()

	req := Request{
		sessionToken: session.Token,
		metrics:      m,
		data:         message,
	}

	var d requestData //figure out how to not crash entire program on bad requests
	if err := json.Unmarshal(message, &d); err != nil {
		fmt.Println(err)
		// return an error, requires some refactoring for the server to handle it
	}

	req.endpoint = d.Endpoint
	req.Params = d.Params

	fmt.Println(d)

	return &req
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

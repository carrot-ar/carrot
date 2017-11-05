package carrot

import (
	//"fmt"
	"time"
)

const (
	endpoint1 = "test1"
	endpoint2 = "test2"
	endpoint3 = "bad_method"
)

type TestController struct {
	count int
}
type TestStreamController struct {
	count int
}

func (c *TestController) Print(req *Request, broadcast *Broadcast) {
	c.count += 1
	//fmt.Printf("This controller's internal count value: %v\n", c.count)
	//broadcast.Send([]byte("This is a controller broadcasting a message!"))

}

func (c *TestStreamController) Print(req *Request, broadcast *Broadcast) {
	c.count += 1
	//fmt.Printf("This stream controller's internal count value: %v\n", c.count)
	//broadcast.Send([]byte("This is a stream controller broadcasting a message!"))
}

func getTokenRouteAndRequestForTest(endpoint string) (SessionToken, *Route, *Request, error) {
	ss := NewDefaultSessionManager()
	token, _, err := ss.NewSession()
	if err != nil {
		return "", nil, nil, err
	}
	sesh, err := ss.Get(token)
	if err != nil {
		return "", nil, nil, err
	}
	req := &Request{
		SessionToken: sesh.Token,
		metrics:      make([]time.Time, MetricCount),
		endpoint:     endpoint,
	}
	route, err := Lookup(req.endpoint)
	if err != nil {
		return "", nil, nil, err
	}
	return token, route, req, nil
}

func getTestController(endpoint string) (*AppController, error) {
	_, route, _, err := getTokenRouteAndRequestForTest(endpoint)
	if err != nil {
		return nil, err
	}
	c, err := NewController(route.Controller(), doCacheControllers)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func init() {
	Environment = "testing"

	Add(endpoint1, TestController{}, "Print", false)
	Add(endpoint2, TestStreamController{}, "Print", true)
	Add(endpoint3, TestController{}, "BadMethod", false)
}

package carrot

import(
	"fmt"
	"time"
)

const (
	endpoint1 = "test1"
	endpoint2 = "test2"
	endpoint3 = "bad_method"
)

type TestDispatcherController struct{}
type TestController struct{}
type TestStreamController struct {
	count int
}

func (tdc *TestDispatcherController) Print(req *Request, broadcast *Broadcast) {
	fmt.Println("The dispatchingRequest method worked because this controller is speaking :D!")
}

func (c *TestController) Print(req *Request, broadcast *Broadcast) {
	fmt.Printf("Hello, world! Here is my event request!!\n")
}

func (c *TestStreamController) Print(req *Request, broadcast *Broadcast) {
	c.count += 1
	fmt.Printf("Stream Controllers internal count value: %v\n", c.count)

	broadcast.Send([]byte("This is the stream controller broadcasting a message!"))
}

func getTokenRouteAndRequestForTest(endpoint string) (SessionToken, *Route, *Request, error) {
	ss := NewDefaultSessionManager()
	token, err := ss.NewSession()
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

	//for dispatcher_test.go
	Add(endpoint1, TestDispatcherController{}, "Print", false)

	//for controller_test.go
	Add(endpoint1, TestController{}, "Print", false)
	Add(endpoint2, TestStreamController{}, "Print", true)
	Add("bad_method", TestController{}, "BadMethod", false)
}

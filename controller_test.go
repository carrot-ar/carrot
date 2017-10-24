package carrot

import (
	"fmt"
	"testing"
	"time"
)

type TestController struct{}
type TestStreamController struct {
	count int
}

func (c *TestController) Print(req *Request, res *Responder) {
	fmt.Printf("Hello, world! Here is my event request!!\n")
	// req.End()
}

func (c *TestStreamController) Print(req *Request, res *Responder) {
	fmt.Printf("Hello, world! Here is my stream request!!\n")
	c.count += 1
	fmt.Printf("Stream Controllers internal count value: %v\n", c.count)
	// req.End()

	res.Broadcast <- []byte("This is the stream controller broadcasting a message!")
}

func TestControllerFactory(t *testing.T) {
	NewController(TestController{}, false)
	NewController(TestStreamController{}, true)
	// handle test
}

func TestMethodInvocation(t *testing.T) {
	// tc := AppController{
	// 	Controller: TestController{},
	// 	persist: false,
	// }
	Add("test1", TestController{}, "Print", false)
	Add("test2", TestStreamController{}, "Print", true)
	// route := Lookup("test")
	// req := NewRequest(nil, nil)

	// tc.Invoke(route, req)
	ss := NewDefaultSessionManager()
	token, err1 := ss.NewSession()
	if err1 != nil {
		fmt.Println(err1)
	}
	sesh, err2 := ss.Get(token)
	if err2 != nil {
		fmt.Println(err2)
	}

	req1 := &Request{
		sessionToken: sesh.Token,
		metrics:      make([]time.Time, MetricCount),
		endpoint:     "test1",
	}

	req2 := &Request{
		sessionToken: sesh.Token,
		metrics:      make([]time.Time, MetricCount),
		endpoint:     "test2",
	}

	req3 := &Request{
		sessionToken: sesh.Token,
		metrics:      make([]time.Time, MetricCount),
		endpoint:     "test2",
	}

	req4 := &Request{
		sessionToken: sesh.Token,
		metrics:      make([]time.Time, MetricCount),
		endpoint:     "test2",
	}

	d := NewDispatcher()
	go d.Run()
	d.requests <- req1
	d.requests <- req2
	d.requests <- req3
	d.requests <- req4
}

func TestInvalidMethodInvocation(t *testing.T) {
	Add("bad_method", TestController{}, "BadMethod", false)
	ss := NewDefaultSessionManager()
	token, err := ss.NewSession()
	if err != nil {
		t.Errorf(err.Error())
	}

	req := &Request{
		sessionToken: token,
		metrics:      make([]time.Time, MetricCount),
		endpoint:     "bad_method",
	}

	rt, _ := Lookup("bad_method")

	c, err := NewController(rt.Controller(), false)
	if err != nil {
		t.Errorf("failed to make controller")
	}

	err = c.Invoke(rt, req)
	if err == nil {
		t.Errorf("Method invocation did not capture invalid method and probably crashed.")
	}

}

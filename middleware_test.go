package buddy

import (
	"testing"
)

var (
	Middleware = NewMiddlewarePipeline()
)

func TestMiddlewareRun(t *testing.T) {
	go Middleware.Run()

	ctx := make(Context)
	ctx["hello"] = "world!"
	Middleware.In <- ctx
	out := <-Middleware.Out

	if ctx["hello"] != out["hello"] {
		t.Errorf("The input context was not equal to the output context: %v != %v", ctx, out)
	}
}

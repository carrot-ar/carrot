package carrot

import (
	"testing"
)

func TestMiddlewareRun(t *testing.T) {
	Middleware := NewMiddlewarePipeline()
	go Middleware.Run()
}

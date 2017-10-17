package buddy

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNewRequest(t *testing.T) {
	sessions := NewDefaultSessionManager()
	token, _ := sessions.NewSession()
	ctx, _ := sessions.Get(token)
	str := `{"endpoint": "test", "params": { "test1": "result1", "test2": "result2" } }`

	actual := NewRequest(ctx, []byte(str))

	expected := &Request{
		sessionToken: ctx.Token,
		endpoint:     "test",
		data:         []byte(str),
		metrics:      make([]time.Time, MetricCount),
	}

	if actual.endpoint != expected.endpoint {
		t.Errorf("actual.endpoint != expected.endpoint... \n %v \n %v \n", actual.endpoint, expected.endpoint)
	}

	if reflect.DeepEqual(actual.Params, expected.Params) {
		t.Errorf("actual != expected... \n %v \n %v \n", actual.Params, expected.Params)
	}
}

func TestRequestTokenMismatch(t *testing.T) {
	sessions := NewDefaultSessionManager()
	token, _ := sessions.NewSession()
	ctx, _ := sessions.Get(token)
	str := `{"session_token": "badsessiontoken", "endpoint": "test", "params": { "test1": "result1", "test2": "result2" } }`

	req := NewRequest(ctx, []byte(str))

	if !strings.Contains(req.err.Error(), "mismatch") {
		t.Errorf("request token validation check failed stop request \n "+
			"Error: %v \n Request Token: %v \n Session Token: %v \n", req.err, req.sessionToken, ctx.Token)
	}
}

package carrot

import (
	"testing"
)

//TODO: test the Transform method and add error handling to helper methods
func TestTransformFunctions(t *testing.T) {
	c := CarrotTransformController{}
	if c.sessions == nil {
		c.sessions = NewDefaultSessionManager()
	}
	primaryToken, err := c.sessions.GetPrimaryDeviceToken()
	if err != nil {
		t.Errorf("There was an error retrieving the primary device token in transform.go")
	}
	_, session, err := c.sessions.NewSession()
	if err != nil {
		t.Errorf("There was an error retrieving the session in transform.go")
	}
	req := &Request{}
	req.Offset, err = NewOffset(3,2,1)
	if (err != nil) {
		t.Error(err)
	}

	c.storeT_L(req, session)

	_, err = c.requestT_P(req)
	if err != nil {
		t.Error(err)
	}	

	c.storeT_P(req, primaryToken)	
}
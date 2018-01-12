package carrot

import (
	// "context"
	"github.com/sirupsen/logrus"
	"encoding/json"
)

type Context interface {
	//context.Context   // still figuring out how we can take advantage of this
	Response() *messageData /* carrot response object */
	Request() *messageData  /* carrot request  object */
	Session() *Session
	Logger() logrus.Entry
	Error() error
	Data() string	// builds response and provides back the json
}


// Carrot Context
type CContext struct {
	request  *messageData
	response *messageData
	session  *Session
	requestParams   Params
	responseParams   Params
	// a new way to log: keep track of state in contexts and use the state as the log field
	logger   *logrus.Entry
	error 	 error

	requestRawMessage []byte
}

func NewCContext(session *Session, message []byte) (*CContext, error) {
	ctx := &CContext{
		response: nil,
		request: nil,
		session: session,
		requestParams: make(map[string]interface{}),
		responseParams: make(map[string]interface{}),
		logger: logrus.WithFields(logrus.Fields{
			"session" : "< setup session token >",
		}),
		error: nil,
		requestRawMessage: message,
	}

	ctx.buildRequest()

	if ctx.Error() != nil {
		return nil, ctx.Error()
	}

	return ctx, nil
}

func (c *CContext) Response() *messageData {
	return c.response
}

func (c *CContext) Request() *messageData {
	return c.request
}

func (c *CContext) Session() *Session {
	return c.session
}

func (c *CContext) Logger() *logrus.Entry {
	return c.logger
}

func (c *CContext) Error() error {
	return c.error
}

func (c *CContext) buildRequest() {

	// Unmarshal json, check for errors and validate, then set the
	// carrot context request to the deserialized message
	var md messageData
	c.error = json.Unmarshal(c.requestRawMessage, &md)

	if c.error != nil {
		return
	}

	c.error = validSession(c.session.Token, SessionToken(md.SessionToken))

	if c.error != nil {
		return
	}

	c.request = &md
}

func (c *CContext) BuildResponse() {
	// call this function to build the response on the delivery of each message to a client
}




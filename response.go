package carrot

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type response interface {
	AddParam()
	AddParams()
}

// for use of controller-level control
type ResponseParams map[string]interface{}

// NewOffset instantiates a struct to hold coordinate data used later to create JSON objects.
func NewOffset(x float64, y float64, z float64) (*offset, error) {
	return &offset{
		X: x,
		Y: y,
		Z: z,
	}, nil
}

// NewPayload instantiates a struct to hold offset and params data used later to create JSON objects.
func NewPayload(sessionToken string, offset *offset, params map[string]interface{}) (payload, error) {
	sessions := NewDefaultSessionManager()
	primaryToken, err := sessions.GetPrimaryDeviceToken()
	if err != nil {
		return payload{}, err
	}

	// if the requester is a primary device,
	// don't do any transform math
	// likewise, if the incoming offset is empty, don't perform a transform
	if sessionToken == string(primaryToken) || offset == nil {
		return newPayloadNoTransform(offset, params)
	}

	currentSession, err := sessions.Get(SessionToken(sessionToken))
	if err != nil {
		return payload{}, err
	}

	//do transform math to get event placement
	e_p, err := getE_P(currentSession, offset)
	//if err != nil {
	//	return payload{}, err
	//}

	log.Info()
	return payload{
		Offset: e_p,
		Params: params,
	}, nil
}

// newPayloadNoTransform instantiates a struct to hold payload data without modifying
// the offset content used later to create JSON objects.
func newPayloadNoTransform(offset *offset, params map[string]interface{}) (payload, error) {
	return payload{
		Offset: offset,
		Params: params,
	}, nil
}

// NewResponse instantiates a struct that will be directly converted to a JSON object to be broadcasted to devices.
func NewResponse(sessionToken string, endpoint string, payload payload) (*messageData, error) {
	return &messageData{
		SessionToken: sessionToken,
		Endpoint:     endpoint,
		Payload:      payload,
	}, nil
}

// AddParam adds a new parameter to the params field of the payload but does not override existing ones.
func (md *messageData) AddParam(key string, value interface{}) {
	if md.Payload.Params == nil {
		md.Payload.Params = make(map[string]interface{})
	}
	params := md.Payload.Params
	_, exists := params[key]
	if !exists {
		params[key] = value
	}
}

// AddParams adds new parameters to the params field of the payload but does not override existing ones.
func (md *messageData) AddParams(rp ResponseParams) {
	for key, value := range rp {
		md.AddParam(key, value)
	}
}

// Build converts response struct into JSON object that will be broadcasted to devices.
func (md *messageData) Build() ([]byte, error) {
	res, err := json.Marshal(md)
	return res, err
}

// CreateDefaultResponse generates a JSON object (ready to broadcast response) from a request in one step.
// This function offers brevity but does not support extra input from controllers (skips adding params).
func CreateDefaultResponse(req *Request) ([]byte, error) {
	payload, err := NewPayload(string(req.SessionToken), req.Offset, req.Params)
	r, err := NewResponse(string(req.SessionToken), req.endpoint, payload)
	res, err := r.Build()
	return res, err
}

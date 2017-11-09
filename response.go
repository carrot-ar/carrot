package carrot

import (
	"encoding/json"
)

type response interface {
	AddParam()
	AddParams()
}

//for use of controller-level control
type ResponseParams map[string]interface{}

func NewOffset(x float64, y float64, z float64) (*offset, error) {
	return &offset {
		X:	x,
		Y:	y,
		Z:	z,
	}, nil
}

func NewPayload(offset *offset, params map[string]string) (payload, error) {
	return payload {
		Offset:	offset,
		Params:	params,
	}, nil
}

func NewResponse(sessionToken string, endpoint string, payload payload) (*messageData, error) {
	return &messageData {
		SessionToken:	sessionToken,
		Endpoint:		endpoint,
		Payload:		payload,
	}, nil
}

//only adds new params, does not override existing ones
func (md *messageData) AddParam(key string, value interface{}) {
	params := md.Payload.Params
	_, exists := params[key]
	if !exists {
		params[key] = value.(string)
	}
}

func (md *messageData) AddParams(rp ResponseParams) {
	for key, value := range rp {
		md.AddParam(key, value)
	}
}

func (md *messageData) Build() ([]byte, error) {
	res, err := json.Marshal(md)
	return res, err
}

//offers brevity but does not support extra input from controllers (skips adding params)
func CreateDefaultResponse(req *Request) ([]byte, error) {
	payload, err := NewPayload(req.Offset, req.Params)
	r, err := NewResponse(string(req.SessionToken), req.endpoint, payload)
	res, err := r.Build()
	return res, err
}

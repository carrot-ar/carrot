package carrot

import (
	"encoding/json"
	//"fmt"
	log "github.com/sirupsen/logrus"
)

type response interface {
	AddParam()
	AddParams()
}

//for use of controller-level control
type ResponseParams map[string]interface{}

func NewOffset(x float64, y float64, z float64) (*offset, error) {
	return &offset{
		X: x,
		Y: y,
		Z: z,
	}, nil
}

func NewPayload(sessionToken string, offset *offset, params map[string]interface{}) (payload, error) {
	sessions := NewDefaultSessionManager()
	primaryToken, err := sessions.GetPrimaryDeviceToken()
	if err != nil {
		return payload{}, err
	}

	// if the requester is a primary device,
	// don't do any transform math
	// likewise, if the incoming offset is empty, don't perform a transform
	if offset == nil { //used to be "if sessionToken == string(primaryToken) || offset == nil {""
		return newPayloadNoTransform(offset, params)
	}

	sender, err := sessions.Get(SessionToken(sessionToken))
	if err != nil {
		return payload{}, err
	}

	//do transform math to get event placement
	//right now only handles transform between single primary and single secondary device
	var recipient *Session
	if sender.isPrimaryDevice() {
		recipient, err = sessions.GetASecondarySession()
		if err != nil {
			return payload{}, err
		}
	} else { //sender is a secondary device
		recipient, err = sessions.Get(SessionToken(primaryToken))
		if err != nil {
			return payload{}, err
		}
	}
	e_p, err := getE_L(sender, recipient, offset)
	if err != nil {
		return payload{}, err
	}
	log.Info()
	return payload{
		Offset: e_p,
		Params: params,
	}, nil
}

func newPayloadNoTransform(offset *offset, params map[string]interface{}) (payload, error) {
	return payload{
		Offset: offset,
		Params: params,
	}, nil
}

func NewResponse(sessionToken string, endpoint string, payload payload) (*messageData, error) {
	return &messageData{
		SessionToken: sessionToken,
		Endpoint:     endpoint,
		Payload:      payload,
	}, nil
}

//only adds new params, does not override existing ones
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
	payload, err := NewPayload(string(req.SessionToken), req.Offset, req.Params)
	r, err := NewResponse(string(req.SessionToken), req.endpoint, payload)
	res, err := r.Build()
	return res, err
}

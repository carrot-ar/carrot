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
	if sessionToken == string(primaryToken) {
		return newPayloadNoTransform(offset, params)
	}

	currentSession, err := sessions.Get(SessionToken(sessionToken))
	if err != nil {
		return payload{}, err
	}

	primaryT_P := currentSession.T_P
	log.Infof("t_p: x: %v y: %v z: %v", primaryT_P.X, primaryT_P.Y, primaryT_P.Z)

	currentT_L := currentSession.T_L

	log.Infof("t_l: x: %v y: %v z: %v", currentT_L.X, currentT_L.Y, currentT_L.Z)

	// offset is the e_l
	// o_p = t_l - t_p
	// e_p = e_l - o_p
	o_p := offsetSub(currentT_L, primaryT_P)
	log.Infof("o_p: x: %v y: %v z: %v", o_p.X, o_p.Y, o_p.Z)

	e_p := offsetSub(offset, o_p)
	log.Infof("e_p: x: %v y: %v z: %v", e_p.X, e_p.Y, e_p.Z)

	log.Infof("e_l: x: %v y: %v z: %v", offset.X, offset.Y, offset.Z)

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

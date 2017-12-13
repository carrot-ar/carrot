package carrot

import (
	"encoding/json"
	//"fmt"
	log "github.com/sirupsen/logrus"
)

type response interface {
	AddParam()
	AddParams()
	build()
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

func NewResponse(sessionToken string, endpoint string, params ResponseParams, request *Request) (*messageData, error) {
	return &messageData{
		SessionToken: sessionToken,
		Endpoint:     endpoint,
		Payload:      nil,
		params:  	  params,
		request: 	  request,
	}, nil
}

//only adds new params, does not override existing ones
func (md *messageData) AddParam(key string, value interface{}) {
	if md.Payload.Params == nil {
		md.Payload.Params = make(ResponseParams)
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

func (md *messageData) build(recipient *Session) ([]byte, error) {
	// prepare the payload for departure
	md.buildPayload(recipient)
	res, err := json.Marshal(md)
	return res, err
}


func (md *messageData) buildPayload (recipient *Session) error {
	sessions := NewDefaultSessionManager() // probably slow
	//primaryToken, err := sessions.GetPrimaryDeviceToken()
	//if err != nil {
	//	return err
	//}

	// if the requester is a primary device,
	// don't do any transform math
	// likewise, if the incoming offset is empty, don't perform a transform
	if md.request.Offset == nil { //used to be "if sessionToken == string(primaryToken) || offset == nil {""
		payload, err := buildPayloadNoTransform(md.request.Offset)
		if  err != nil {
			return err
		}

		md.Payload = payload
	}

	sender, err := sessions.Get(SessionToken(md.SessionToken))
	if err != nil {
		return err
	}

	/*

	//do transform math to get event placement
	var recipient *Session
	if sender.isPrimaryDevice() {
		//recipient, err = sessions.GetASecondarySession()
		//if err != nil {
		//	return err
		//}
		recipient = secondarySession
	} else { //sender is a secondary device
		recipient, err = sessions.Get(SessionToken(primaryToken))
		if err != nil {
			return err
		}
	}

	*/

	e_p, err := getE_P(sender, recipient, md.request.Offset)
	if err != nil {
		return err
	}
	log.Info()

	md.Payload = payload{
		Offset: e_p,
		Params: md.params,
	}

	return nil
}


func buildPayloadNoTransform(offset *offset) (payload, error) {
	return payload{
		Offset: offset,
		Params: nil,
	}, nil
}


/*
// NOTE: ignoring for now for prototyping  purposes
//offers brevity but does not support extra input from controllers (skips adding params)
func CreateDefaultResponse(req *Request) ([]byte, error) {
	payload, err := NewPayload(string(req.SessionToken), req.Offset, req.Params)
	r, err := NewResponse(string(req.SessionToken), req.endpoint, payload)
	res, err := r.Build()
	return res, err
}
*/

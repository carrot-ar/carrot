package carrot

import (
	"encoding/json"
)

func NewOffset(x float64, y float64, z float64) (*offset, error) {
	return &offset{
		X: x,
		Y: y,
		Z: z,
	}, nil
}

func (md *messageData) build(recipient *Session) ([]byte, error) {
	// prepare the payload for departure
	md.buildPayload(recipient)
	res, err := json.Marshal(md)
	return res, err
}

func (md *messageData) buildPayload(recipient *Session) error {
	//_ := NewDefaultSessionManager()

	/*
		Killian's logic
	*/

	/*
	md.Payload = payload{
		Offset: nil, // calculated offset
		Params: md.params,
	}
	*/
	return nil
}

func buildPayloadNoTransform(offset *offset) (payload, error) {
	return payload{
		Offset: offset,
		Params: nil,
	}, nil
}

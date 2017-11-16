package carrot

/*
	{
		"session_token": "KjIQhKUPNrvHkUHv1VySBg==",
		"endpoint": "test_endpoint",
		"payload": {
			"offset": {
				"x": 3.2,
				"y": 1.3,
				"z": 4.0
			},
			"params": {
				"foo": "bar"
			}
		}
	}
*/

// represents incoming requests and outgoing responses
type messageData struct {
	SessionToken string  `json:"session_token"`
	Endpoint     string  `json:"endpoint"`
	Payload      payload `json:"payload"`
}

type payload struct {
	Offset *offset                `json:"offset,omitempty"`
	Params map[string]interface{} `json:"params,omitempty"`
}

type offset struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

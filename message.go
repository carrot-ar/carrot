package buddy

/*
	{
		"session_token": "KjIQhKUPNrvHkUHv1VySBg==",
		"endpoint": "test_endpoint",
		"origin": {
			"longitude": 45.501689,
			"latitude": -73.567256
		},
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

// incoming message
type requestData struct {
	SessionToken string      `json:"session_token"`
	Endpoint     string      `json:"endpoint"`
	Origin       originData  `json:"origin"`
	Payload      payloadData `json:"payload"`
}

type originData struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type payloadData struct {
	Offset offsetData        `json:"offset"`
	Params map[string]string `json:"params"`
}

type offsetData struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}
package carrot

/*

This file describes the construction and types of messages sent during the Picnic Protocol handshake
between the server and primary and secondary devices.


 Example initial message to devices. Since the session_token and uuid are the same,
 this is a primary device message. The two fields would be different for secondary devices.

{
	session_token: "E621E1F8-C36C-495A-93FC-0C247A3E6E5F";
	endpoint: "carrot_beacon";
	payload: {
	  params: {
			identifier = "com.Carrot.PrimaryBeacon";
			uuid = "E621E1F8-C36C-495A-93FC-0C247A3E6E5F";
	  };
	};
  }

*/

func createInitialDeviceInfo(uuid string, token string) ([]byte, error) {
	params := ResponseParams{"identifier": "com.Carrot.Beacon", "uuid": uuid}
	payload, err := newPayloadNoTransform(nil, params)
	res, err := NewResponse(token, "carrot_beacon", payload)
	info, err := res.Build()
	return info, err
}

/*

 Example request to obtain T_P from the primary device.
 Since params (and offsets) are unnecessary for this case, the payload is empty.

{
	session_token: "E621E1F8-C36C-495A-93FC-0C247A3E6E5F";
	endpoint: "carrot_transform";
	payload: {

	};
}

*/

func getT_PFromPrimaryDeviceRes(token string) ([]byte, error) {
	payload, err := newPayloadNoTransform(nil, nil)
	res, err := NewResponse(token, "carrot_transform", payload)
	ask, err := res.Build()
	return ask, err
}

/*

Example response from a primary device sending its T_P or a secondary device sending its T_L.
These messages look identical so the context in transform.go determines where and what the offset is stored as.

{
	"session_token": "E621E1F8-C36C-495A-93FC-0C247A3E6E5F",
	"endpoint": "carrot_transform",
	"payload": {
		"offset": {
			"x": 1,
			"y": 1,
			"z": 1
		}
	}
}

*/

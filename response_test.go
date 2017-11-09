package carrot

import (
	"testing"
	//"fmt"
	"encoding/json"
)

func TestBuildResponse(t *testing.T) {
	token := "totally_a_real_token"
	endpoint := "test_endpoint"
	x, y, z := 3.2, 1.3, 4.0
	offset, err := NewOffset(x, y, z)
	if err != nil {
		t.Error(err)
	}
	params := map[string]string{"red": "fish", "blue": "fish"}
	payload_complete, err := NewPayload(offset, params)
	if err != nil {
		t.Error(err)
	}
	payload_noparams, err := NewPayload(offset, nil)
	if err != nil {
		t.Error(err)
	}
	payload_nooffset, err := NewPayload(nil, params)
	if err != nil {
		t.Error(err)
	}
	payload_empty, err := NewPayload(nil, nil)
	if err != nil {
		t.Error(err)
	}
	r_c, err := NewResponse(token, endpoint, payload_complete)
	if err != nil {
		t.Error(err)
	}
	r_np, err := NewResponse(token, endpoint, payload_noparams)
	if err != nil {
		t.Error(err)
	}
	r_no, err := NewResponse(token, endpoint, payload_nooffset)
	if err != nil {
		t.Error(err)
	}
	r_e, err := NewResponse(token, endpoint, payload_empty)
	if err != nil {
		t.Error(err)
	}
	res_c, err := r_c.Build()
	if err != nil {
		t.Error(err)
	}
	res_np, err := r_np.Build()
	if err != nil {
		t.Error(err)
	}
	res_no, err := r_no.Build()
	if err != nil {
		t.Error(err)
	}
	res_e, err := r_e.Build()
	if err != nil {
		t.Error(err)
	}
	//fmt.Printf("%s\n", res_c)
	//fmt.Printf("%s\n", res_np)
	//fmt.Printf("%s\n", res_no)
	//fmt.Printf("%s\n", res_e)
	if !isJSON(res_c)  {
		t.Error("The complete response is not valid JSON")
	}
	if !isJSON(res_np)  {
		t.Error("The response with no params is not valid JSON")
	}
	if !isJSON(res_no)  {
		t.Error("The response with no offset is not valid JSON")
	}
	if !isJSON(res_e)  {
		t.Error("The response with an empty payload is not valid JSON")
	} 

	// r_c.AddParam("newparam", "test")
	// r_np.AddParam("newparam", "test")
	// r_no.AddParam("newparam", "test")
	// r_e.AddParam("newparam", "test")
	// res_c, err = r_c.Build()
	// if err != nil {
	// 	t.Error(err)
	// }
	// res_np, err = r_np.Build()
	// if err != nil {
	// 	t.Error(err)
	// }
	// res_no, err = r_no.Build()
	// if err != nil {
	// 	t.Error(err)
	// }
	// res_e, err = r_e.Build()
	// if err != nil {
	// 	t.Error(err)
	// }
	// if !isJSON(res_c)  {
	// 	t.Error("The complete response with an extra added param is not valid JSON")
	// }
	// if !isJSON(res_np)  {
	// 	t.Error("The response with a single added param is not valid JSON")
	// }
	// if !isJSON(res_no)  {
	// 	t.Error("The response with no offset and an extra added param is not valid JSON")
	// }
	// if !isJSON(res_e)  {
	// 	t.Error("The response with an added param and otherwise empty payload is not valid JSON")
	// } 

	//TODO: test all variations of CreateDefaultResponse
	s := NewDefaultSessionManager()
	sessionToken, _, err := s.NewSession()
	req_c := newResponseTestRequest(sessionToken, endpoint, params, offset)
	res, err := CreateDefaultResponse(req_c)
	if !isJSON(res)  {
		t.Error("The default response is not valid JSON")
	}
	req_np := newResponseTestRequest(sessionToken, endpoint, nil, offset)
	res, err = CreateDefaultResponse(req_np)
	if !isJSON(res)  {
		t.Error("The default response is not valid JSON")
	}
	req_no := newResponseTestRequest(sessionToken, endpoint, params, nil)
	res, err = CreateDefaultResponse(req_no)
	if !isJSON(res)  {
		t.Error("The default response is not valid JSON")
	} 
	req_e := newResponseTestRequest(sessionToken, endpoint, nil, nil)
	res, err = CreateDefaultResponse(req_e)
	if !isJSON(res)  {
		t.Error("The default response is not valid JSON")
	} 
}

func isJSON(res []byte) bool {
    var js json.RawMessage
    return json.Unmarshal([]byte(res), &js) == nil
}

func newResponseTestRequest(t SessionToken, e string, p map[string]string, o *offset) *Request {
	return &Request {
		SessionToken: t,
		endpoint:	  e,
		Params:		  p,
		Offset:	      o,
		metrics:      nil,
		data:         nil,
	}
}
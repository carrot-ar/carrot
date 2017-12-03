package carrot

import (
	"testing"
	//"fmt"
	"encoding/json"
)

func TestBuildResponse(t *testing.T) {
	sm := NewDefaultSessionManager()
	_, primarySession, err := sm.NewSession()
	if err != nil {
		t.Error(err)
	}
	primarySession.primaryDevice = true
	token, secondarySession, err := sm.NewSession()
	if err != nil {
		t.Error(err)
	}
	secondarySession.T_L, err = NewOffset(3, 2, 1)
	if err != nil {
		t.Error(err)
	}
	secondarySession.T_P, err = NewOffset(6, 5, 4)
	if err != nil {
		t.Error(err)
	}
	endpoint := "test_endpoint"
	x, y, z := 3.2, 1.3, 4.0
	offset, err := NewOffset(x, y, z)
	if err != nil {
		t.Error(err)
	}

	params := make(map[string]interface{})
	params["one"] = "fish"
	params["two"] = "fish"

	//test building a response from scratch
	payload_complete, err := NewPayload(string(token), offset, params)
	if err != nil {
		t.Error(err)
	}
	payload_noparams, err := NewPayload(string(token), offset, nil)
	if err != nil {
		t.Error(err)
	}
	payload_nooffset, err := NewPayload(string(token), nil, params)
	if err != nil {
		t.Error(err)
	}
	payload_empty, err := NewPayload(string(token), nil, nil)
	if err != nil {
		t.Error(err)
	}
	r_c, err := NewResponse(string(token), endpoint, payload_complete)
	if err != nil {
		t.Error(err)
	}
	r_np, err := NewResponse(string(token), endpoint, payload_noparams)
	if err != nil {
		t.Error(err)
	}
	r_no, err := NewResponse(string(token), endpoint, payload_nooffset)
	if err != nil {
		t.Error(err)
	}
	r_e, err := NewResponse(string(token), endpoint, payload_empty)
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
	if !isJSON(res_c) {
		t.Error("The complete response is not valid JSON")
	}
	if !isJSON(res_np) {
		t.Error("The response with no params is not valid JSON")
	}
	if !isJSON(res_no) {
		t.Error("The response with no offset is not valid JSON")
	}
	if !isJSON(res_e) {
		t.Error("The response with an empty payload is not valid JSON")
	}

	//test AddParam and AddParams
	testKey, testValue := "newparam", "test"
	testMap1 := ResponseParams{"red": "fish", "blue": "fish"}
	testMap2 := ResponseParams{"blue": "fish", "new": "fish"}
	r_c.AddParam(testKey, testValue)
	_, ok := r_c.Payload.Params[testKey]
	if !ok {
		t.Error("AddParam did not work correctly on a complete response")
	}
	r_np.AddParam(testKey, testValue)
	_, ok = r_np.Payload.Params[testKey]
	if !ok {
		t.Error("AddParam did not work correctly on a response without parameters")
	}
	r_no.AddParam(testKey, testValue)
	_, ok = r_no.Payload.Params[testKey]
	if !ok {
		t.Error("AddParam did not work correctly on a response without an offset")
	}
	r_e.AddParam(testKey, testValue)
	_, ok = r_e.Payload.Params[testKey]
	if !ok {
		t.Error("AddParam did not work correctly on an empty response")
	}
	r_c.AddParams(testMap1)
	_, ok = r_c.Payload.Params["red"]
	if !ok {
		t.Error("AddParams did not work correctly on a complete response")
	}
	r_c.AddParams(testMap2)
	_, ok = r_c.Payload.Params["new"]
	if !ok {
		t.Error("AddParams ignored a value it was supposed to add")
	}
	_, ok = r_c.Payload.Params["fake"]
	if ok {
		t.Error("AddParams tried to return a value that wasn't assigned")
	}
	//TODO: test more edge cases of AddParams
	res_c, err = r_c.Build()
	if err != nil {
		t.Error(err)
	}
	res_np, err = r_np.Build()
	if err != nil {
		t.Error(err)
	}
	res_no, err = r_no.Build()
	if err != nil {
		t.Error(err)
	}
	res_e, err = r_e.Build()
	if err != nil {
		t.Error(err)
	}
	if !isJSON(res_c) {
		t.Error("The complete response with an extra added param is not valid JSON")
	}
	if !isJSON(res_np) {
		t.Error("The response with a single added param is not valid JSON")
	}
	if !isJSON(res_no) {
		t.Error("The response with no offset and an extra added param is not valid JSON")
	}
	if !isJSON(res_e) {
		t.Error("The response with an added param and otherwise empty payload is not valid JSON")
	}

	//test all variations of CreateDefaultResponse
	s := NewDefaultSessionManager()
	sessionToken, _, err := s.NewSession()
	req_c := newResponseTestRequest(sessionToken, endpoint, params, offset)
	res, err := CreateDefaultResponse(req_c)
	if !isJSON(res) {
		t.Error("The default response is not valid JSON")
	}
	req_np := newResponseTestRequest(sessionToken, endpoint, nil, offset)
	res, err = CreateDefaultResponse(req_np)
	if !isJSON(res) {
		t.Error("The default response is not valid JSON")
	}
	req_no := newResponseTestRequest(sessionToken, endpoint, params, nil)
	res, err = CreateDefaultResponse(req_no)
	if !isJSON(res) {
		t.Error("The default response is not valid JSON")
	}
	req_e := newResponseTestRequest(sessionToken, endpoint, nil, nil)
	res, err = CreateDefaultResponse(req_e)
	if !isJSON(res) {
		t.Error("The default response is not valid JSON")
	}
}

func isJSON(res []byte) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(res), &js) == nil
}

func newResponseTestRequest(t SessionToken, e string, p map[string]interface{}, o *offset) *Request {
	return &Request{
		SessionToken: t,
		endpoint:     e,
		Params:       p,
		Offset:       o,
		metrics:      nil,
		data:         nil,
	}
}

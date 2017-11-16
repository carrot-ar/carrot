package carrot

import (
	"testing"
)

func TestCreateInitialDeviceInfo(t *testing.T) {
	//TODO: create test with hardcoded string comparison
	uuid, token := "someUUID", "testToken"
	_, err := createInitialDeviceInfo(uuid, token)
	if err != nil {
		t.Error(err)
	}
}

func TestGetT_PFromPrimaryDeviceRes(t *testing.T) {
	//TODO: create test with hardcoded string comparison
	token := "testToken"
	_, err := getT_PFromPrimaryDeviceRes(token)
	if err != nil {
		t.Error(err)
	}
}
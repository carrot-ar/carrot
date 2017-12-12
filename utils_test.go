package carrot

import "testing"

func TestInSlice(t *testing.T) {
	items := []string{"potato", "radish", "carrot"}
	existsStr := "potato"
	doesNotExistStr := "beet"

	if !InSlice(existsStr, items) {
		t.Errorf("%v was in items when it shouldn't have been", existsStr)
	}

	if InSlice(doesNotExistStr, items) {
		t.Errorf("%v was in items when it shouldn't have been", doesNotExistStr)
	}
}

func TestGenerateUUID(t *testing.T) {
	_, err := generateUUID()
	if err != nil {
		t.Error(err)
	}
}

func TestOffsetSub(t *testing.T) {
	a := &offset{
		X: 1,
		Y: 1,
		Z: 1,
	}

	b := &offset{
		X: 1,
		Y: 1,
		Z: 1,
	}

	c := offsetSub(a, b)
	if c.X != 0 || c.Y != 0 || c.Z != 0 {
		t.Errorf("all generic calculations should equal zero: X: %v, Y: %v, Z: %v", c.X, c.Y, c.Z)
	}
}

func TestGetE_P(t *testing.T) {
	store := NewDefaultSessionManager()
	_, session, err := store.NewSession()
	if err != nil {
		t.Error(err)
	}
	primaryToken, err := store.GetPrimaryDeviceToken()
	if err != nil {
		t.Error(err)
	}
	primarySession, err := store.Get(SessionToken(primaryToken))
	if err != nil {
		t.Error(err)
	}
	//make sure you are testing two different session roles
	if session == primarySession {
		_, session, err = store.NewSession()
	}
	session.T_L = &offset{
		X: 2,
		Y: 2,
		Z: 2,
	}
	session.T_P = &offset{
		X: 1,
		Y: 1,
		Z: 1,
	}
	testOffset := &offset{
		X: 1,
		Y: 1,
		Z: 1,
	}
	//test secondary sender and primary recipient
	e_p, err := getE_P(session, primarySession, testOffset)
	if err != nil {
		t.Error(err)
	}
	if e_p.X != 0 || e_p.Y != 0 || e_p.Z != 0 {
		t.Errorf("all E_P values should equal zero: X: %v, Y: %v, Z: %v", e_p.X, e_p.Y, e_p.Z)
	}
	//test primary sender and secondary recipient
	e_p, err = getE_P(primarySession, session, testOffset)
	if err != nil {
		t.Error(err)
	}
	if e_p.X != 0 || e_p.Y != 0 || e_p.Z != 0 {
		t.Errorf("all E_P values should equal zero: X: %v, Y: %v, Z: %v", e_p.X, e_p.Y, e_p.Z)
	}
}

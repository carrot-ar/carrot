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

	b:= &offset{
		X: 1,
		Y: 1,
		Z: 1,
	}

	c := offsetSub(a, b)
	if c.X != 0 || c.Y != 0 || c.Z != 0 {
		t.Errorf("should all equal zero: X: %v, Y: %v, Z: %v", c.X, c.Y, c.Z)
	}
}

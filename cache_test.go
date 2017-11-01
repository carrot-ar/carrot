package carrot

import (
	"testing"
)

func TestCacheInstantiation(t *testing.T) {
	ccl := NewCachedControllersList()
	if ccl == nil {
		t.Error("The cached controllers list was not instantiated")
	}
}

func TestCacheGetAndAdd(t *testing.T) {
	key := "key"
	ccl := NewCachedControllersList()

	if ccl.Exists(key) {
		t.Error("The key shouldn't exist in the cached controller list yet")
	}

	_, err := ccl.Get(key)
	if err == nil {
		t.Error("A controller shouldn't have been returned because its key doesn't exist")
	}

	c, err := getTestController(endpoint1)
	if err != nil {
		t.Error(err)
	}
	ccl.Add(key, c)

	if !ccl.Exists(key) {
		t.Error("The key should exist in the cached controllers list")
	}

	_, err = ccl.Get(key)
	if err != nil {
		t.Error(err)
	}
}	

func TestCacheDeleteOldest(t *testing.T) {
	key := "key"
	ccl := NewCachedControllersList()

	err := ccl.DeleteOldest()
	if err == nil {
		t.Error("Nothing should have been deleted since none of the cached controllers have a matching key")
	}

	c, err := getTestController(endpoint1)
	if err != nil {
		t.Error(err)
	}
	ccl.Add(key, c)

	err = ccl.DeleteOldest()
	if err != nil {
		t.Error(err)
	}
}
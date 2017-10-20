package carrot

import (
	"reflect"
)

// returns a new controller of controllerType
func NewController(c interface{}, isStream bool) (*AppController, error) {

	newController := reflect.New(reflect.TypeOf(c))
	responder := NewResponder()
	go responder.Run()

	return &AppController{
		Controller: newController,
		persist:    isStream,
		responder:  responder,
	}, nil
}

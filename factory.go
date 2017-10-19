package carrot

import (
	"reflect"
)

// returns a new controller of controllerType
func NewController(c interface{}, isStream bool) (*AppController, error) {

	newController := reflect.New(reflect.TypeOf(c))

	if isStream {
		return &AppController{
			Controller: newController,
			persist:    true,
		}, nil
	}
	return &AppController{
		Controller: newController,
		persist:    false,
	}, nil
}

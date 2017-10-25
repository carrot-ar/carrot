package carrot

import (
	"reflect"
)

// returns a new controller of controllerType
func NewController(c interface{}, isStream bool) (*AppController, error) {

	newController := reflect.New(reflect.TypeOf(c))


	return &AppController{
		Controller: newController,
		persist:    isStream,
		broadcast:  broadcast,
	}, nil
}

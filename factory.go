package carrot

import (
	"reflect"
)

// returns a new controller of controllerType
func NewController(c interface{}, isStream bool) (*AppController, error) {

	newController := reflect.New(reflect.TypeOf(c))

	broadcast := NewBroadcast(broadcaster)

	go broadcast.broadcaster.Run()

	return &AppController{
		Controller: newController,
		persist:    isStream,
		broadcast:  broadcast,
	}, nil
}

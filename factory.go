package carrot

import (
	"reflect"
)

// NewController initializes a new controller instance of controllerType.
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

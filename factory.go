package buddy

import (
	"reflect"
	"fmt"
)

const controllerInitializer = "Initialize"

// returns a new controller of controllerType
func NewController(c interface{}, isStream bool) (*AppController, error) {

	fmt.Println("MAKING NEW CONTROLLER")
	newController := reflect.New(reflect.TypeOf(c))
	fmt.Println(reflect.TypeOf(c).Name())
	initializer := newController.MethodByName(controllerInitializer)

	if initializer.IsValid() {
		args := []reflect.Value{}
		initializer.Call(args)
	} else {
		return nil, fmt.Errorf("error: invalid controller")
	}

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

package buddy

import (
	"fmt"
)



// Base controller
type ControllerType string

type DefaultController struct {}


// returns a new controller of controllerType
func New(t ControllerType) (*AppController, error) {
	switch t {
	case "DefaultController":
		return &AppController{
			Controller: DefaultController{},
			persist: false,
		}, nil
	default:
		return nil, fmt.Errorf("error: invalid controller type %v", t)
	}
}






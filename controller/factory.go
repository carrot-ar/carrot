package controller

import (
	"fmt"
)

const (
	/*
		Add controllers below this line
	 */
	DEFAULT_CONTROLLER = iota
)


// Base controller
type ControllerType int

type DefaultController struct {}


// returns a new controller of controllerType
func New(t ControllerType) (*AppController, error) {
	switch t {
	case DEFAULT_CONTROLLER:
		return &AppController{
			Controller: DefaultController{},
			persist: false,
		}, nil
	default:
		return nil, fmt.Errorf("error: invalid controller type %v", t)
	}
}


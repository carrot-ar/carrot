package buddy

import "fmt"

// returns a new controller of controllerType
func NewController(c interface{}) (*AppController, error) {
	fmt.Printf("We're about to create this controller! %v\n", c)
	return &AppController{
		Controller: c,
		persist: false,
	}, nil
}






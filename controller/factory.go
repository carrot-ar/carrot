package controller

import( "github.com/senior-buddy/buddy"

)

const (
	/*
		Add controllers below this line
	 */
	appController = iota
)

type ControllerType int

type Controller interface{
	Invoke(string, *buddy.Request) error
	Persist(bool)
}

/*
	AppController. This is a default controller on which simple routes can be made.
	Adding a route such as "place_ball => AppController_PlaceBall" will route to the
	PlaceBall function within the controller
 */
type AppController struct {
	persist bool
}

func (c *AppController) Persist(p bool) {
	c.persist = p
}

func (c *AppController) Invoke(function string, request *buddy.Request) error {
	// TODO: make this call functions
	return nil
}

// returns a new controller of controllerType
func New(t ControllerType) (Controller, error) {
	switch t {
	case appController:
		return &AppController{ persist: false }, nil
	default:
		return &AppController{ persist: false }, nil
	}
}

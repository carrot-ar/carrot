package buddy

// returns a new controller of controllerType
func NewController(c interface{}) (*AppController, error) {
	return &AppController{
		Controller: c,
		persist: false,
	}, nil
}






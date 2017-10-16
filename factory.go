package buddy

// returns a new controller of controllerType
func NewController(c interface{}, isStream bool) (*AppController, error) {
	if isStream {
		return &AppController{
			Controller: c,
			persist:    true,
		}, nil
	}
	return &AppController{
		Controller: c,
		persist:    false,
	}, nil
}

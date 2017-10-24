package carrot

type Route struct {
	controller ControllerType
	function   string
	persist    bool
}

type Endpoint string

func (r *Route) Controller() ControllerType {
	return r.controller
}

func (r *Route) Function() string {
	return r.function
}

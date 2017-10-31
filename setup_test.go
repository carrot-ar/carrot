package carrot

const (
	//for dispatcher_test.go
	endpoint1 = "test1"
)

func init() {
	Environment = "testing"	

	//for dispatcher_test.go
	Add(endpoint1, TestDispatcherController{}, "Print", false)	
}

package controller

import (
	"testing"
	"fmt"
	"github.com/senior-buddy/buddy"
	"github.com/senior-buddy/buddy/routes"
	"reflect"
)

func TestControllerFactory(t *testing.T) {
	//c, _ := New(appController)
	//fmt.Println(c)
	//fmt.Println(reflect.TypeOf(c))
}
/*
type TestController struct {
	base Controller
}

func (c *TestController) Print(req *buddy.Request) {
	fmt.Printf("Hello, world!\n")
}

func TestMethodInvocation(t *testing.T) {
	tc := TestController{}

	b := Controller{
		persist: false,
		parent: &tc,
	}

	tc.base = b


	route := routes.Lookup("test")
	req := buddy.Request{}

	tc.base.Invoke(route, &req)
}
*/

type TestController struct {}

func (c *TestController) Print(req *buddy.Request) {
	fmt.Printf("Hello, world! Here is my request!! %v\n", req)
}

func (c *TestController) PrintMore(req *buddy.Request) {
	fmt.Printf("I am working, yay! Here is my request!! %v\n", req)
}

func TestMethodInvocation(t *testing.T) {
	tc := BaseController{
		Controller: TestController{},
		persist: false,
	}

	fmt.Println(reflect.TypeOf(tc.Controller))



	for {
		route := routes.Lookup("test")
		req := buddy.NewRequest(nil, nil)

		tc.Invoke(route, req)
	}



}
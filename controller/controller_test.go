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

func (c *DefaultController) Print(req *buddy.Request) {
	fmt.Printf("Hello, world! Here is my request!!\n")
	req.End()
}

func TestMethodInvocation(t *testing.T) {
	tc := AppController{
		Controller: DefaultController{},
		persist: false,
	}

	fmt.Println(reflect.TypeOf(tc.Controller))

	route := routes.Lookup("test")
	req := buddy.NewRequest(nil, nil)

	tc.Invoke(route, req)

}
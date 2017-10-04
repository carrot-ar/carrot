package controller

import (
	"testing"
	"fmt"
)

func TestControllerFactory(t *testing.T) {
	c, _ := New(appController)
	fmt.Println(c)
	//fmt.Println(reflect.TypeOf(c))
}

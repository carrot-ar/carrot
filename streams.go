package buddy

import (
	"fmt"
)

type OpenStreamsList struct {
	streams map[SessionToken]*AppController
}

func NewOpenStreamsList() *OpenStreamsList {
	return &OpenStreamsList{
		streams: make(map[SessionToken]*AppController),
	}
}

func (osl *OpenStreamsList) Exists(token SessionToken) bool {
	_, ok := osl.streams[token]
	if ok {
		return true
	}
	return false
}

func (osl *OpenStreamsList) Get(token SessionToken) *AppController {
	sc, ok := osl.streams[token]
	if !ok || sc == nil {
		fmt.Println("cannot return route because it doesn't exist")
		//return nil
	}
	fmt.Printf("%v:%v\n", token, &sc)
	return sc
}

func (osl *OpenStreamsList) Add(token SessionToken, ac *AppController) {
	fmt.Printf("%v:%v\n", token, &ac)
	osl.streams[token] = ac
}

func (osl *OpenStreamsList) Delete(token SessionToken) {
	delete(osl.streams, token) //doesn't return anything
}

func (osl *OpenStreamsList) IsEmpty() bool {
	if len(osl.streams) == 0 {
		return true
	}
	return false
}

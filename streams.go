package buddy

import (
	//nothing yet
)

type OpenStreamsList struct {
	streams map[SessionToken]*AppController
	length int
}

func NewOpenStreamsList() *OpenStreamsList {
	return &OpenStreamsList {
		streams:	make(map[SessionToken]*AppController),
		length:		0,
	}
}

func (osl *OpenStreamsList) Exists(token SessionToken) bool {
	_, ok := osl.streams[token]
	if ok { return true }
	return false
}

func (osl *OpenStreamsList) Get(token SessionToken) *AppController {
	sc, ok := osl.streams[token]
	if !ok {
		println("cannot return route be it doesn't exist")
		return nil
	}
	return sc
}

func (osl *OpenStreamsList) Add(token SessionToken) error {
	//osl.streams[token] = NewController()
	osl.length += 1
	return nil //get rid of this
}

func (osl *OpenStreamsList) Delete(token SessionToken) {
	delete(osl.streams, token) //doesn't return anything
	osl.length -= 1
}

func (osl *OpenStreamsList) IsEmpty() bool {
	if osl.length == 0 { return true }
	return false
}

package main

import(
	"github.com/senior-buddy/carrot"
)

func main() {

	/*
		init sequence
	 */

	// start the session manager here

	// start the server here

	// do this to be able to handle large client counts
	// ulimit -n SOME_REALLY_BIG_NUMBER
	//
	server := carrot.NewServer()
	go server.Middleware.Run()
	go server.Run()
	server.Serve()
}

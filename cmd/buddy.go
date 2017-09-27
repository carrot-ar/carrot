package main

import(
	"github.com/senior-buddy/buddy"
)

func main() {

	/*
		init sequence
	 */

	// start the session manager here

	// start the server here
	server := buddy.NewServer()
	go server.Middleware.Run()
	go server.Run()
	server.Serve()
}

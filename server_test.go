package buddy

import (
	"flag"
	"log"
	"net/http"
	"testing"
)

var (
	addr = flag.String("addr", ":8080", "http service address")
)

func TestToSeeIfTheServerRuns(t *testing.T) {
	flag.Parse()
	server := NewServer()
	go server.Run()

	//this block needs to be updated for buddy-esque purposes
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(server, w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		t.Error("something did not work and you should be sad for this reason")
		t.Errorf("The error is: %d", err)
	}
}

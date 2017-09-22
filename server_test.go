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

//this method may be unnecessary later on
func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func TestToSeeIfTheServerRuns(t *testing.T) {
	flag.Parse()
	server := newServer()
	go server.run()

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

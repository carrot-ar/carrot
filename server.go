package carrot

import (
	"flag"
	"fmt"
	"net/http"
	log "github.com/sirupsen/logrus"
)

const (
	serverSecret         = "37FUqWlvJhRgwPMM1mlHOGyPNwkVna3b"
	broadcastChannelSize = 4096
	port                 = 8080
)

//the server maintains the list of clients and
//broadcasts messages to the clients
type Server struct {

	//register requests from the clients
	register chan *Client

	//unregister requests from the clients
	unregister chan *Client

	//access list of existing sessions
	sessions SessionStore

	//keep track of middleware
	Middleware *MiddlewarePipeline
}

func NewServer(sessionStore SessionStore) *Server {
	return &Server{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		sessions:   sessionStore,
		Middleware: NewMiddlewarePipeline(),
	}
}

func (svr *Server) Run() {
	for {
		select {
		case client := <-svr.register:
			client.softOpen()
			token := <-client.sendToken
			//create persistent token for new or invalid sessions
			exists := svr.sessions.Exists(token)
			if (token == "") || !exists {
				var err error
				token, err = svr.sessions.NewSession()
				if err != nil {
					//handle later
					log.Print(err)
				}
				//return the new token for the session
				client.sendToken <- token
			}

			svr.sessions.SetClient(token, client)
			close(client.start)
		case client := <-svr.unregister:
			if client.Open() {
				client.softClose()
				// delete client?
				close(client.send)
				close(client.sendToken)
				client = nil
			}
		}
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	//log.Println(r.URL)

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

func (svr *Server) Serve() {
	addr := flag.String("addr", fmt.Sprintf(":%d", port), "http service address")

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(svr, w, r)
	})

	log.WithFields(log.Fields{
		"port": port,
		"url": "ws://localhost/",
	}).Infof("Listening...")

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		fmt.Println(err)
	}
}

package carrot

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

const (
	serverSecret         = "37FUqWlvJhRgwPMM1mlHOGyPNwkVna3b"
	broadcastChannelSize = 65536
	port                 = 8080
	maxClients           = 4096
)

type Clients [maxClients]*Client

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

	clients *Clients
}

func NewServer(sessionStore SessionStore) *Server {
	return &Server{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		sessions:   sessionStore,
		Middleware: NewMiddlewarePipeline(),
		clients: 	new(Clients),
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
				token, sessionPtr, err := svr.sessions.NewSession()
				if err != nil {
					//handle later
					log.Print(err)
				}

				client.session = sessionPtr

				// find a free location in the client list
				for i, c := range svr.clients {
					if c == nil {
						svr.clients[i] = client
						break
					}
				}

				//return the new token for the session
				client.sendToken <- token
			}

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

	log.Printf("Listening at http://localhost:%d", port)
	log.Printf("Listening at ws://localhost:%d", port)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		fmt.Println(err)
	}
}

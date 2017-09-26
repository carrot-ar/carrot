package buddy

import (
	"net/http"
	"flag"
	"log"
	"fmt"
)

const serverSecret = "37FUqWlvJhRgwPMM1mlHOGyPNwkVna3b"

//the server maintains the list of clients and
//broadcasts messages to the clients
type Server struct {

	//registered clients
	clients map[*Client]bool

	//inbound messages from the clients
	broadcast chan []byte

	//register requests from the clients
	register chan *Client

	//unregister requests from the clients
	unregister chan *Client

	//access list of existing sessions
	sessions SessionStore
}

func NewServer() *Server {
	return &Server{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		sessions:   NewDefaultSessionManager(),
	}
}

func (svr *Server) Run() {
	for {
		select {
		case client := <-svr.register:
			svr.clients[client] = true
			token := <-client.sendToken
			//create persistent token for new or invalid sessions
			exists := svr.sessions.Exists(token)
			if (token == "nil") || !exists {
				var err error
				token, err = svr.sessions.NewSession()
				if err != nil {
					//handle later
				}
				//return the new token for the session
				client.sendToken <- token
			}
		case client := <-svr.unregister:
			if _, ok := svr.clients[client]; ok {
				//delete corresponding session
				delete(svr.clients, client)
				close(client.send)
				close(client.sendToken)
			}
		case message := <-svr.broadcast:
			for client := range svr.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					close(client.sendToken)
					delete(svr.clients, client)
				}
			}
		}
	}
}

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
	http.ServeFile(w, r, "../home.html")
}

func (svr *Server) Serve() {
	addr := flag.String("addr", ":8080", "http service address")

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(svr, w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		fmt.Println(err)
	}
}
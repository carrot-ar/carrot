package buddy

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

const serverSecret = "37FUqWlvJhRgwPMM1mlHOGyPNwkVna3b"

//the server maintains the list of clients and
//broadcasts messages to the clients
type Server struct {

	//registered clients
	//clients map[*Client]bool

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
		//clients:    make(map[*Client]bool),
		sessions: NewDefaultSessionManager(),
	}
}

func (svr *Server) Run() {
	for {
		select {
		case client := <-svr.register:
			client.open = true
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

			svr.sessions.SetClient(token, client)

		case client := <-svr.unregister:
			if client.open {
				// This is leaving the context in existence which may eventually
				// fill up, so either we need to delete this context on close (hard)
				// or use this boolean as a mark and sweep tactic like with gc (easy)
				client.open = false
				//delete corresponding session
				//delete(svr.clients, client)
				close(client.send)
				close(client.sendToken)
				client = nil
			}
		case message := <-svr.broadcast:
			svr.broadcastAll(message)
		}
	}
}

func (svr *Server) broadcastAll(message []byte) {
	start := time.Now()
	svr.sessions.Range(func(key, value interface{}) bool {
		ctx := value.(*Context)

		if ctx.SessionExpired() {
			svr.sessions.Delete(ctx.Token)
			return true
		} else if !ctx.Client.open {
			return true
		}

		select {
		case ctx.Client.send <- message:
			return true
		default:
			close(ctx.Client.send)
			close(ctx.Client.sendToken)
		}

		return false
	})
	end := time.Now()
	fmt.Printf("Time to broadcast to %v users: %v\n",
		svr.sessions.Length(),
		end.Sub(start))
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
	http.ServeFile(w, r, "home.html")

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

package buddy

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

const (
	serverSecret         = "37FUqWlvJhRgwPMM1mlHOGyPNwkVna3b"
	broadcastChannelSize = 4096
	port                 = 8080
)

//the server maintains the list of clients and
//broadcasts messages to the clients
type Server struct {
	//inbound messages from the clients
	broadcast chan []byte

	//register requests from the clients
	register chan *Client

	//unregister requests from the clients
	unregister chan *Client

	//access list of existing sessions
	sessions SessionStore

	//keep track of middleware
	Middleware *MiddlewarePipeline
}

func NewServer() *Server {
	return &Server{
		broadcast:  make(chan []byte, broadcastChannelSize),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		sessions:   NewDefaultSessionManager(),
		Middleware: NewMiddlewarePipeline(),
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
			if client.open {
				client.open = false
				// delete client?
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
	// start := time.Now()
	expiredSessionCount := 0
	closedClientCount := 0
	refreshedClientCount := 0
	messagesSent := 0
	svr.sessions.Range(func(key, value interface{}) bool {
		ctx := value.(*Session)

		if ctx.SessionExpired() {
			expiredSessionCount++
			svr.sessions.Delete(ctx.Token)
			return true
		} else if !ctx.Client.open {
			closedClientCount++
			return true
		}

		ctx.expireTime = refreshExpiryTime()
		refreshedClientCount++

		select {
		case ctx.Client.send <- message:
			messagesSent++
			return true
		default:
			close(ctx.Client.send)
			close(ctx.Client.sendToken)
		}

		return false
	})
	log.Printf("server: broadcast sent %v, refresh %v, closed %v, expired %v",
		messagesSent,
		refreshedClientCount,
		closedClientCount,
		expiredSessionCount)
	// end := time.Now()
	// fmt.Printf("Time to broadcast to %v users: %v\n",
	// 	svr.sessions.Length(),
	// 	end.Sub(start))
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

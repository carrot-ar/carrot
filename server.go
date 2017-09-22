package buddy

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

func newServer() *Server {
	return &Server{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		sessions:   NewDefaultSessionManager(),
	}
}

func (svr *Server) run() {
	for {
		select {
		case client := <-svr.register:
			svr.clients[client] = true
			token := <-client.sendToken
			//create persistent token for new sessions
			if token == "nil" {
				var err error
				token, err = svr.sessions.NewSession()
				if err != nil {
					//handle later
				}
			}
			//make sure that a token exists for the session
			client.sendToken <- token
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

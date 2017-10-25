package carrot

type Broadcaster struct {
	sessions SessionStore
	clients *Clients
	//inbound messages from the clients
	broadcast chan []byte
}

func NewBroadcaster(clients *Clients) *Broadcaster {
	return &Broadcaster{
		sessions:  NewDefaultSessionManager(),
		broadcast: make(chan []byte, broadcastChannelSize),
		clients: clients,
	}
}

func (br *Broadcaster) broadcastAll(message []byte) {
	//log.Printf("Length of clients: %v\n", len(br.clients))
	for _, client := range br.clients {
		if client != nil {
			if client.Expired() {
				br.sessions.Delete(client.session.Token)
				continue
			} else if !client.Open() {
				continue
			}

			client.session.expireTime = refreshExpiryTime()

			client.send <- message
		}
	}
}

func (br *Broadcaster) Run() {
	for {
		select {
		case message := <-br.broadcast:
			br.broadcastAll(message)
		}
	}
}

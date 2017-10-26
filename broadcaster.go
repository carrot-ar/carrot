package carrot

type OutboundMessage struct {
	message []byte
}

// manage broadcast groups with the broadcaster
type Broadcaster struct {
	sessions SessionStore
	clientPool  Pool
	//inbound messages from the clients
	broadcast chan []byte
}

func NewBroadcaster(pool Pool) *Broadcaster {
	return &Broadcaster{
		sessions:  NewDefaultSessionManager(),
		broadcast: make(chan []byte, broadcastChannelSize),
		clientPool: pool,
	}
}

// use functions are arguments for broadcasting to the correct groups
func (br *Broadcaster) broadcastAll(message []byte) {
	outboundMessage := OutboundMessage{
		// set the criteria function here
		message: message,
	}
	br.clientPool.Send(&outboundMessage)
}

func (br *Broadcaster) Run() {
	for {
		select {
		case message := <-br.broadcast:
			br.broadcastAll(message)
		}
	}
}

package carrot

type Responder struct {
	sessions SessionStore

	//inbound messages from the clients
	Broadcast chan []byte
}

func NewResponder() *Responder {
	return &Responder{
		sessions:  NewDefaultSessionManager(),
		Broadcast: make(chan []byte, broadcastChannelSize),
	}
}

func (res *Responder) broadcastAll(message []byte) {
	//log.Printf("Current length of this responder: %v\n", len(res.Broadcast))
	messagesSent := 0
	res.sessions.Range(func(key, value interface{}) bool {
		ctx := value.(*Session)


		if ctx.SessionExpired() {
			res.sessions.Delete(ctx.Token)
			return true
		} else if !ctx.Client.Open() {
			return true
		}

		ctx.expireTime = refreshExpiryTime()

		select {
		case ctx.Client.send <- message:
			messagesSent++
			return true
		}

		return true
	})
	//log.Printf("server: broadcast sent %v, refresh %v, closed %v, expired %v",
	//	messagesSent,
	//	refreshedClientCount,
	//	closedClientCount,
	//	expiredSessionCount)
}

func (res *Responder) Run() {
	for {
		select {
		case message := <-res.Broadcast:
			go res.broadcastAll(message)
			//fmt.Println(string(message))
		}
	}
}

package carrot

var broadcast *Broadcast

func Run() {
	sessions := NewDefaultSessionManager()
	server := NewServer(sessions)
	dispatcher := NewDispatcher()
	broadcaster := NewBroadcaster(server.clients)
	broadcast = NewBroadcast(broadcaster)
	go broadcast.broadcaster.Run()
	go dispatcher.Run()
	go server.Middleware.Run()
	go server.Run()
	server.Serve()
}

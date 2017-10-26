package carrot

var broadcast *Broadcast

func Run() {
	sessions := NewDefaultSessionManager()
	server := NewServer(sessions)
	dispatcher := NewDispatcher()
	clientPool := NewClientPool()
	broadcaster := NewBroadcaster(clientPool)
	broadcast = NewBroadcast(broadcaster)
	go clientPool.ListenAndSend()
	go broadcast.broadcaster.Run()
	go dispatcher.Run()
	go server.Middleware.Run()
	go server.Run()
	server.Serve()
}

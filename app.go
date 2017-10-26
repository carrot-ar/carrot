package carrot

var broadcast *Broadcast

func Run() {
	sessions := NewDefaultSessionManager()
	clientPool := NewClientPool()
	server := NewServer(clientPool, sessions)
	dispatcher := NewDispatcher()
	broadcaster := NewBroadcaster(clientPool)
	broadcast = NewBroadcast(broadcaster)
	go broadcast.broadcaster.Run()
	go dispatcher.Run()
	go server.Middleware.Run()
	go server.Run()
	go broadcast.broadcaster.clientPool.ListenAndSend()
	server.Serve()
}

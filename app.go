package carrot

var config = getConfig() // parse conf.yaml and instantiate Config

func Run() {
	sessions := NewDefaultSessionManager()
	server := NewServer(sessions)
	dispatcher := NewDispatcher()
	go dispatcher.Run()
	go server.Middleware.Run()
	go server.Run()
	server.Serve()
}

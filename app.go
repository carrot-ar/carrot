package carrot

var config Config

func Run() {
	config = getConfig() // parse conf.yaml and instantiate Config
	sessions := NewDefaultSessionManager()
	server := NewServer(sessions)
	dispatcher := NewDispatcher()
	go dispatcher.Run()
	go server.Middleware.Run()
	go server.Run()
	server.Serve()
}

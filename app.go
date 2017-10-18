package buddy

func Run() {
	sessions := NewDefaultSessionManager()
	server := NewServer(sessions)
	go server.Middleware.Run()
	go server.Run()
	server.Serve()
}

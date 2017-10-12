package buddy

func Run() {
	server := NewServer()
	go server.Middleware.Run()
	go server.Run()
	server.Serve()
}

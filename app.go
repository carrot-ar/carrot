package buddy

func Run(rt *RoutingTable) {
	server := NewServer()
	go server.Middleware.Run()
	go server.Run()
	server.Serve()
}

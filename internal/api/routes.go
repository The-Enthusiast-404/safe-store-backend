package api

func (s *Server) SetupRoutes() {
	s.router.GET("/", s.HelloWorld)
}

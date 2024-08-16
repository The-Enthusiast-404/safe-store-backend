package api

func (s *Server) SetupRoutes() {
	s.router.GET("/", s.HelloWorld)
	s.router.POST("/upload", s.UploadFile)
	s.router.GET("/download/:filename", s.DownloadFile)
}

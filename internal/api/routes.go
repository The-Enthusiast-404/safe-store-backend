package api

import (
	"dev.theenthusiast.safe-store/internal/middleware"
)

func (s *Server) SetupRoutes() {
	s.router.GET("/", middleware.CORS(s.HelloWorld))
	s.router.POST("/upload", middleware.CORS(s.UploadFile))
	s.router.GET("/download/:filename", middleware.CORS(s.DownloadFile))
	s.router.GET("/files", middleware.CORS(s.ListFiles))
}

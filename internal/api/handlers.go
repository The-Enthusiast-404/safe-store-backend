package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) HelloWorld(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s.logger.Info("Handling Hello World request")
	fmt.Fprint(w, "Hello, World!")
}

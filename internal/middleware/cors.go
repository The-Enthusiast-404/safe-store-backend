package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func CORS(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Disposition, Content-Length")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition, Content-Length")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r, ps)
	}
}

func HandleCORS() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Disposition, Content-Length")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition, Content-Length")
		w.WriteHeader(http.StatusOK)
	})
}

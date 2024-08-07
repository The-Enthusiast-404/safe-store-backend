package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func (app *application) routes() http.Handler {
	// initialize a new httprouter instance
	router := httprouter.New()

	// Convert the notFoundResponse() helper to a http.Handler using the
	// http.HandlerFunc() adapter, and then set it as the custom error handler for 404
	// Not Found responses.
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// Convert the methodNotAllowedResponse() helper to a http.Handler using the http.HandlerFunc() adapter, and then set it as the custom error handler for 405 Method Not Allowed responses.
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// register the healthcheck handler function with the router
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/upload", app.uploadFileHandler)
	router.HandlerFunc(http.MethodGet, "/v1/download/:filename", app.downloadFileHandler)
	router.HandlerFunc(http.MethodGet, "/v1/files", app.listFilesHandler) // New endpoint
	router.HandlerFunc(http.MethodPost, "/v1/delete", app.deleteFilesHandler)

	// Create a new CORS handler
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Add your frontend origin here
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
	})

	return c.Handler(router)
}

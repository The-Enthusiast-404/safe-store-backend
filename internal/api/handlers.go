package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/julienschmidt/httprouter"
)

func (s *Server) HelloWorld(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s.logger.Info("Handling Hello World request")
	fmt.Fprint(w, "Hello, World!")
}

func (s *Server) UploadFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = s.r2Client.UploadFile(r.Context(), header.Filename, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File %s uploaded successfully", header.Filename)
}

func (s *Server) DownloadFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filename := ps.ByName("filename")

	result, err := s.r2Client.DownloadFile(r.Context(), filename)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer result.Body.Close()

	// Set the correct headers
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", aws.ToString(result.ContentType))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", aws.ToInt64(result.ContentLength)))

	// Stream the file content to the response writer
	_, err = io.Copy(w, result.Body)
	if err != nil {
		http.Error(w, "Failed to send file", http.StatusInternalServerError)
		return
	}
}

func (s *Server) ListFiles(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	files, err := s.r2Client.ListFiles(r.Context())
	if err != nil {
		http.Error(w, "Failed to list files", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

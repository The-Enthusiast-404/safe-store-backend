package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/julienschmidt/httprouter"
)

func (app *application) uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "File too large")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid file")
		return
	}
	defer file.Close()

	// Upload the file to Digital Ocean Spaces
	_, err = app.s3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(app.config.doSpacesBucket),
		Key:    aws.String(header.Filename),
		Body:   file,
		ACL:    aws.String("private"),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				app.logger.Printf("Bucket does not exist: %v", aerr)
				app.serverErrorResponse(w, r, fmt.Errorf("storage is not properly configured"))
			case "AccessDenied":
				app.logger.Printf("Access denied: %v", aerr)
				app.errorResponse(w, r, http.StatusForbidden, "Access denied")
			default:
				app.logger.Printf("S3 error: %v", aerr)
				app.serverErrorResponse(w, r, err)
			}
		} else {
			app.logger.Printf("Unknown error: %v", err)
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"message": "File uploaded successfully"}, nil)
}

func (app *application) downloadFileHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	filename := params.ByName("filename")

	// Get the file from Digital Ocean Spaces
	result, err := app.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(app.config.doSpacesBucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	defer result.Body.Close()

	// Set the appropriate headers
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", *result.ContentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", *result.ContentLength))

	// Stream the file to the response
	_, err = io.Copy(w, result.Body)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) listFilesHandler(w http.ResponseWriter, r *http.Request) {
	// List objects in the bucket
	resp, err := app.s3.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(app.config.doSpacesBucket),
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Prepare the file list
	var files []map[string]interface{}
	for _, item := range resp.Contents {
		files = append(files, map[string]interface{}{
			"name": *item.Key,
			"size": *item.Size,
		})
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"files": files}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteFileHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	filename := params.ByName("filename")

	// Delete the file from Digital Ocean Spaces
	_, err := app.s3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(app.config.doSpacesBucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				app.notFoundResponse(w, r)
			case "AccessDenied":
				app.logger.Printf("Access denied: %v", aerr)
				app.errorResponse(w, r, http.StatusForbidden, "Access denied")
			default:
				app.logger.Printf("S3 error: %v", aerr)
				app.serverErrorResponse(w, r, err)
			}
		} else {
			app.logger.Printf("Unknown error: %v", err)
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "File deleted successfully"}, nil)
}

func (app *application) deleteFilesHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Filenames []string `json:"filenames"`
	}

	err := app.readJSON(w, r, &req)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if len(req.Filenames) == 0 {
		app.errorResponse(w, r, http.StatusBadRequest, "No files specified for deletion")
		return
	}

	var objects []*s3.ObjectIdentifier
	for _, filename := range req.Filenames {
		objects = append(objects, &s3.ObjectIdentifier{Key: aws.String(filename)})
	}

	deleteInput := &s3.DeleteObjectsInput{
		Bucket: aws.String(app.config.doSpacesBucket),
		Delete: &s3.Delete{
			Objects: objects,
			Quiet:   aws.Bool(false),
		},
	}

	result, err := app.s3.DeleteObjects(deleteInput)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Check for any errors during deletion
	if len(result.Errors) > 0 {
		errorMessages := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			errorMessages[i] = fmt.Sprintf("Failed to delete %s: %s", *err.Key, *err.Message)
		}
		app.errorResponse(w, r, http.StatusInternalServerError, errorMessages)
		return
	}

	deletedCount := len(result.Deleted)
	message := fmt.Sprintf("Successfully deleted %d file(s)", deletedCount)
	app.writeJSON(w, http.StatusOK, envelope{"message": message}, nil)
}

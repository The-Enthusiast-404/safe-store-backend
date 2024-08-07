package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Config struct will be used to hold all the configuration settings of the application
type config struct {
	port           int
	env            string
	doSpacesKey    string
	doSpacesSecret string
	doSpacesRegion string
	doSpacesBucket string
}

// Application struct will be used to hold all the dependencies of the application
type application struct {
	config config
	logger *log.Logger
	s3     *s3.S3
}

func main() {
	// instance of config struct
	var cfg config

	// reading configuration settings from environment variables
	cfg.port = 4000 // You can make this configurable
	cfg.env = os.Getenv("APP_ENV")
	cfg.doSpacesKey = os.Getenv("DO_SPACES_KEY")
	cfg.doSpacesSecret = os.Getenv("DO_SPACES_SECRET")
	cfg.doSpacesRegion = os.Getenv("DO_SPACES_REGION")
	cfg.doSpacesBucket = os.Getenv("DO_SPACES_BUCKET")

	// instance of logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Create a new AWS session
	s3Session, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(cfg.doSpacesKey, cfg.doSpacesSecret, ""),
		Endpoint:         aws.String(fmt.Sprintf("https://%s.digitaloceanspaces.com", cfg.doSpacesRegion)),
		Region:           aws.String(cfg.doSpacesRegion),
		S3ForcePathStyle: aws.Bool(false),
	})
	if err != nil {
		logger.Fatal(err)
	}

	// Create a new S3 client
	s3Client := s3.New(s3Session)

	app := &application{
		config: cfg,
		logger: logger,
		s3:     s3Client,
	}

	// Check if the bucket exists, if not, create it
	err = app.ensureBucketExists()
	if err != nil {
		logger.Fatal(err)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// starting the server
	logger.Printf("Starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

func (app *application) ensureBucketExists() error {
	_, err := app.s3.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(app.config.doSpacesBucket),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound":
				// Bucket doesn't exist, create it
				_, err = app.s3.CreateBucket(&s3.CreateBucketInput{
					Bucket: aws.String(app.config.doSpacesBucket),
				})
				if err != nil {
					return fmt.Errorf("failed to create bucket: %v", err)
				}
				app.logger.Printf("Bucket %s created successfully", app.config.doSpacesBucket)
			case "Forbidden":
				return fmt.Errorf("forbidden: unable to access bucket %s. Check your permissions", app.config.doSpacesBucket)
			default:
				return fmt.Errorf("error checking bucket %s: %v", app.config.doSpacesBucket, err)
			}
		}
	}

	return nil
}

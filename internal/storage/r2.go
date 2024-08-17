package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"dev.theenthusiast.safe-store/internal/config"
)

type R2Client struct {
	client *s3.Client
	bucket string
}

func NewR2Client(cfg *config.Config) (*R2Client, error) {
	r2Resolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.R2AccountID),
		}, nil
	})

	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithEndpointResolver(r2Resolver),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.R2AccessKeyID, cfg.R2AccessKeySecret, "")),
		awsconfig.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg)

	return &R2Client{
		client: client,
		bucket: cfg.R2Bucket,
	}, nil
}

func (r *R2Client) UploadFile(ctx context.Context, key string, body io.Reader) error {
	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
		Body:   body,
	})
	return err
}

func (r *R2Client) DownloadFile(ctx context.Context, key string) (*s3.GetObjectOutput, error) {
	result, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

type File struct {
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
}

func (r *R2Client) ListFiles(ctx context.Context) ([]File, error) {
	result, err := r.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(r.bucket),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	var files []File
	for _, obj := range result.Contents {
		file := File{
			Name: aws.ToString(obj.Key),
			Size: aws.ToInt64(obj.Size),
		}
		if obj.LastModified != nil {
			file.LastModified = *obj.LastModified
		}
		files = append(files, file)
	}

	return files, nil
}

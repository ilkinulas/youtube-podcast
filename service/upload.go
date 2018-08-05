package service

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ilkinulas/youtube-podcast/config"
)

type Uploader interface {
	Upload(filename string) (string, error)
}

type S3Uploader struct {
	logger *log.Logger
	cfg    config.S3
}

func NewS3Uploader(cfg config.S3, logger *log.Logger) Uploader {
	return &S3Uploader{
		logger: logger,
		cfg:    cfg,
	}
}

func (u *S3Uploader) Upload(filename string) (string, error) {
	u.logger.Printf("Uploading file %q to bucket %q", filename, u.cfg.Bucket)
	awsCfg := aws.NewConfig().
		WithRegion(u.cfg.Regioin).
		WithMaxRetries(10).
		WithS3ForcePathStyle(true)
	if u.cfg.Endpoint != "" {
		awsCfg = awsCfg.WithEndpoint(u.cfg.Endpoint)
	}
	if u.cfg.Key != "" {
		awsCfg = awsCfg.WithCredentials(
			credentials.NewStaticCredentials(u.cfg.Key, u.cfg.Secret, ""))
	}
	s3Session, err := session.NewSessionWithOptions(
		session.Options{
			Config:            *awsCfg,
			SharedConfigState: session.SharedConfigEnable,
		})
	if err != nil {
		return "", fmt.Errorf("failed to create s3 session %v", err)
	}

	uploader := s3manager.NewUploader(s3Session)

	f, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to open file %q, %v", filename, err)
	}

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(u.cfg.Bucket),
		Key:    aws.String(filename), //TODO change
		Body:   f,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %v\n", result.Location)
	return result.Location, nil
}

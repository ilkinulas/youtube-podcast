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
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ilkinulas/youtube-podcast/storage"
	"net/url"
)

type Uploader interface {
	Upload(video storage.Video) (string, error)
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

func (u *S3Uploader) Upload(video storage.Video) (string, error) {
	u.logger.Printf("Uploading file %q to bucket %q", video.Filename, u.cfg.Bucket)
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

	f, err := os.Open(video.Filename)
	if err != nil {
		return "", fmt.Errorf("failed to open file %q, %v", video.Filename, err)
	}

	key, err := u.generateS3Key(video)
	if err != nil {
		return "", err
	}
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(u.cfg.Bucket),
		Key:    aws.String(key),
		Body:   f,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file, %v", err)
	}
	u.logger.Printf("file uploaded to, %v\n", result.Location)

	return u.getPresignedUrl(s3Session, key)
}

func (u *S3Uploader) generateS3Key(video storage.Video) (string, error) {
	pUrl, err := url.Parse(video.YoutubeUrl)
	if err != nil {
		return "", err
	}

	values, err := url.ParseQuery(pUrl.RawQuery)
	if err != nil {
		return "", err
	}
	if id, ok := values["v"]; ok {
		return id[0], nil
	}
	return "", fmt.Errorf("failed to extract video id from url %q", video.YoutubeUrl)
}

func (u *S3Uploader) getPresignedUrl(session *session.Session, key string) (string, error) {
	awsS3 := s3.New(session)
	req, _ := awsS3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(u.cfg.Bucket),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(u.cfg.PresignedUrlDuration.Duration)

	if err != nil {
		return "", err
	}

	return urlStr, nil
}

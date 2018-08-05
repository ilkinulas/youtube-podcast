package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ilkinulas/youtube-podcast/storage"
)

type Service struct {
	storage    storage.Storage
	logger     *log.Logger
	ctx        context.Context
	downloader Downloader
	uploader   Uploader
}

func NewService(
	ctx context.Context,
	storage storage.Storage,
	downloader Downloader,
	uploader Uploader,
	logger *log.Logger,
) *Service {
	return &Service{
		ctx:        ctx,
		storage:    storage,
		logger:     logger,
		downloader: downloader,
		uploader:   uploader,
	}
}

func (s *Service) Loop() {
	for {
		select {
		case <-time.After(1 * time.Second):
			err := s.handleNextUrl()
			if err != nil {
				s.logger.Printf("Failed to handle url, %v.", err)
				continue
			}
		case <-s.ctx.Done():
			s.logger.Printf("Download service closing.")
			return
		}
	}
}

func (s *Service) handleNextUrl() error {
	url, err := s.storage.NextUrl()
	if err != nil {
		return err
	}
	if url == "" {
		return nil
	}
	video, err := s.downloadWithRetries(url, 3)
	if err != nil {
		storageErr := s.storage.DownloadFailed(url)
		if storageErr != nil {
			return fmt.Errorf("failed to mark url as failed. url %v, %v", url, err)
		}
		return fmt.Errorf("failed to download url %v. %v", url, err)
	}
	err = s.storage.Downloaded(url)
	if err != nil {
		return fmt.Errorf("failed to mark url %v as downloaded. %v", url, err)
	}
	s.logger.Printf("Video Downloaded ! %v", video.Title)
	// upload video
	uploadUrl, err := s.uploader.Upload(video.Filename)
	if err != nil {
		s.logger.Printf("Failed to upload file %q, %v", video.Filename, err)
	}

	s.logger.Printf("Upload url : %v", uploadUrl)
	// update storage

	return nil
}

func (s *Service) downloadWithRetries(url string, numTries int) (*storage.Video, error) {
	var (
		attempt = 0
		err     error
		video   *storage.Video
	)
	for attempt < numTries {
		attempt++
		video, err = s.downloader.Download(url)
		if err != nil {
			<-time.After(3 * time.Second)
			continue
		}
		return video, nil
	}
	return nil, err
}

func (s *Service) execWithRetry(f func() error, numTries int) error {
	attempt := 0
	var err error
	for attempt < numTries {
		attempt++
		err = f()
		if err != nil {
			<-time.After(3 * time.Second)
			continue
		}
		return nil
	}
	return err
}

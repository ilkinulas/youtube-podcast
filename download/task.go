package download

import (
	"github.com/ilkinulas/youtube-podcast/storage"
	"log"
	"context"
	"time"
	"fmt"
	"os/exec"
	"os"
)

type Service struct {
	storage storage.Storage
	logger  *log.Logger
	ctx     context.Context
}

func NewService(ctx context.Context, storage storage.Storage, logger *log.Logger) *Service {
	return &Service{
		ctx:     ctx,
		storage: storage,
		logger:  logger,
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
	err = s.download(url)
	if err != nil {
		return fmt.Errorf("failed to download url %v. %v", url, err)
	}
	err = s.storage.Downloaded(url)
	if err != nil {
		return fmt.Errorf("failed to mark url %v as downloaded. %v", url, err)
	}
	return nil
}

func (s *Service) download(url string) error {
	s.logger.Printf("Downloading url %v", url)

	cmd := exec.Command("youtube_download.py", url)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

package service

import (
	"log"
	"os/exec"

	"github.com/ilkinulas/youtube-podcast/storage"
	"os"
)

type Downloader interface {
	Download(url string) (*storage.Video, error)
}

func NewPythonDownloader(logger *log.Logger) Downloader {
	return &PythonDownloader{
		logger: logger,
	}
}

type PythonDownloader struct {
	logger *log.Logger
}

func (d *PythonDownloader) Download(url string) (*storage.Video, error) {
	command := exec.Command("./youtube_download.py", url)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	out, err := command.Output()

	if err != nil {
		return nil, err
	}
	video, err := storage.NewVideo(string(out[:]))
	if err != nil {
		return nil, err
	}
	return &video, nil
}

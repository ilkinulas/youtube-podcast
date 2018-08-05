package download

import (
	"github.com/ilkinulas/youtube-podcast/storage"
	"os/exec"
	"log"
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
	out, err := exec.Command("./youtube_download.py", url).Output()
	if err != nil {
		return nil, err
	}
	video, err := storage.NewVideo(string(out[:]))
	if err != nil {
		return nil, err
	}
	d.logger.Printf("Video info downloaded for  %v", video.Title)
	cmd := exec.Command("./youtube_download.py", url, "-d")
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return &video, nil
}

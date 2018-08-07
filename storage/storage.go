package storage

type Storage interface {
	Add(url string) error
	Downloaded(url string) error
	DownloadFailed(url string) error
	NextUrl() (string, error)
	SaveVideo(v Video) error
	SelectVideos() ([]Video, error)
}

const (
	UrlStatusNew            = 0
	UrlStatusDownloaded     = 1
	UrlStatusDownloadFailed = 2
)

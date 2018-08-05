package storage

type Storage interface {
	Add(url string) error
	Downloaded(url string) error
	DownloadFailed(url string) error
	NextUrl() (string, error)
	SaveVideo(v Video) error
}

const (
	URL_STATUS_NEW             = 0
	URL_STATUS_DOWNLOADED      = 1
	URL_STATUS_DOWNLOAD_FAILED = 2
)

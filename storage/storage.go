package storage

type Storage interface {
	Add(url string) error
	Downloaded(url string) error
	NextUrl() (string, error)
}

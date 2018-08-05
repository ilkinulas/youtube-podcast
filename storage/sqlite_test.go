package storage

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NextUrl_Should_Return_Error_If_Table__Is_Empty(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	storage, err := NewSqliteStorage(f.Name())

	assert.NoError(t, err)

	url, err := storage.NextUrl()
	assert.Equal(t, "", url)
	assert.NoError(t, err)
}

func Test_NextUrl_Should_Return_Oldest_Url(t *testing.T) {
	f, _ := ioutil.TempFile("", "")
	storage, _ := NewSqliteStorage(f.Name())

	assert.NoError(t, storage.Add("http://youtube.com/1"))
	assert.NoError(t, storage.Add("http://youtube.com/2"))
	assert.NoError(t, storage.Add("http://youtube.com/3"))

	url, err := storage.NextUrl()
	assert.NoError(t, err)
	assert.Equal(t, "http://youtube.com/1", url)

	assert.NoError(t, storage.Downloaded("http://youtube.com/1"))
	url, err = storage.NextUrl()
	assert.NoError(t, err)
	assert.Equal(t, "http://youtube.com/2", url)
}

func Test_Url_Must_Be_Unique(t *testing.T) {
	f, _ := ioutil.TempFile("", "")
	storage, _ := NewSqliteStorage(f.Name())

	url := "http://youtube.com/1"
	assert.NoError(t, storage.Add(url))
	assert.Error(t, storage.Add(url))
}

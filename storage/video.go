package storage

import "encoding/json"

type Video struct {
	Filename  string `json:filename`
	Url       string `json:url`
	Title     string `json:title`
	Length    int    `json:length`
	Thumbnail string `json:thumb`
}

func NewVideo(s string) (Video, error) {
	var video Video
	err := json.Unmarshal([]byte(s), &video)
	return video, err
}

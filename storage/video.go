package storage

import (
	"encoding/json"
	"fmt"
)

type Video struct {
	Filename   string `json:"filename"`
	YoutubeUrl string `json:"url"`
	Title      string `json:"title"`
	Length     int    `json:"length"`
	Thumbnail  string `json:"thumb"`
	Author     string `json:"author"`
	PublicUrl  string
}

func NewVideo(s string) (Video, error) {
	var video Video
	err := json.Unmarshal([]byte(s), &video)
	if err != nil {
		fmt.Printf("Faied to unmarshal json %v", err)
	}
	return video, err
}

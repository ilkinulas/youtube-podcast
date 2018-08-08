package httpserver

import (
	"fmt"
	"net/http"

	"github.com/ilkinulas/youtube-podcast/storage"
	"github.com/gorilla/feeds"
	"log"
	"github.com/ilkinulas/youtube-podcast/version"
	"time"
)

func SaveUrl(storage storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urls, ok := r.URL.Query()["url"]
		if !ok || len(urls[0]) < 1 {
			http.Error(w, "missing url parameter", http.StatusBadRequest)
			return
		}
		url := urls[0]
		if err := storage.Add(url); err != nil {
			http.Error(w, fmt.Sprintf("failed to save url %v", url), http.StatusInternalServerError)
			return
		}
		println(url)
		w.Write([]byte("OK"))
	})
}

func Rss(log *log.Logger, storage storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		videos, err := storage.SelectVideos()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch videos from storage, %v", err), http.StatusInternalServerError)
			return
		}
		feed := &feeds.Feed{
			Title: "Ilkin's youtube podcast",
			Link: &feeds.Link{
				Href: "https://ilkinulas.github.io",
			},
		}

		var feedItems []*feeds.Item
		for _, video := range videos {
			item := &feeds.Item{
				Title:       video.Title,
				Description: video.Thumbnail,
				Link: &feeds.Link{
					Href: video.PublicUrl,
				},
				Id: "test1",
				Created: time.Now(),
				Updated: time.Now(),
				Enclosure: &feeds.Enclosure{
					Length: fmt.Sprintf("%v", video.Length),
					Url:    video.PublicUrl,
					Type:   "video/mpeg",
				},
			}
			feedItems = append(feedItems, item)
		}
		if len(feedItems) > 0 {
			feed.Items = feedItems
		}
		rss, err := feed.ToRss()
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte(rss))
	})
}

func Index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		versionStr := fmt.Sprintf("Youtube Podcast App %v", version.GetHumanVersion())
		w.Write([]byte(versionStr))
	})
}

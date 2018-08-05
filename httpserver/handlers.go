package httpserver

import (
	"fmt"
	"net/http"

	"github.com/ilkinulas/youtube-podcast/storage"
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

func Index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Index"))
	})
}

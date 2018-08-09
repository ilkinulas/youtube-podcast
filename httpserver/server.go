package httpserver

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ilkinulas/youtube-podcast/config"
	"github.com/ilkinulas/youtube-podcast/storage"
)

type Server struct {
	ctx        context.Context
	logger     *log.Logger
	httpServer *http.Server
	storage    storage.Storage
}

func NewServer(ctx context.Context, listenAddr string, logger *log.Logger, storage storage.Storage, cfg config.Config) *Server {
	router := http.NewServeMux()
	router.Handle("/", Index())
	router.Handle("/save", SaveUrl(storage))
	router.Handle("/rss", Rss(logger, storage, cfg.Podcast))

	server := &http.Server{
		Addr:         listenAddr,
		Handler:      router,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		ctx:        ctx,
		httpServer: server,
		logger:     logger,
		storage:    storage,
	}
}

func (s *Server) Start() error {
	s.logger.Printf("Starting http server on %v", s.httpServer.Addr)
	go func() {
		<-s.ctx.Done()
		s.logger.Printf("Shutting down http server.")
		s.httpServer.Shutdown(s.ctx)
	}()

	return s.httpServer.ListenAndServe()
}

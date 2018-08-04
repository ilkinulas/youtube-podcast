package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"context"
	"net/http"
	"github.com/ilkinulas/youtube-podcast/httpserver"
	"flag"
	"github.com/ilkinulas/youtube-podcast/version"
	"github.com/ilkinulas/youtube-podcast/storage"
	"github.com/ilkinulas/youtube-podcast/download"
)

var (
	listenAddr   string
	loggerPrefix string
	dbFile       string
)

func main() {
	parseFlags()

	logger := log.New(os.Stdout, loggerPrefix, log.LstdFlags)
	logger.Printf("Starting youtube-podcast app %v", version.GetHumanVersion())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
		<-sigch
		cancel()
	}()

	storage, err := storage.NewSqliteStorage(dbFile)
	if err != nil {
		logger.Fatalf("Failed to initialize storage, %v", err)
	}

	downloadService := download.NewService(ctx, storage, logger)
	go downloadService.Loop()

	server := httpserver.NewServer(ctx, listenAddr, logger, storage)
	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Failed to start http server at %v, %v", listenAddr, err)
	}
}

func parseFlags() {
	flag.StringVar(&listenAddr, "listen-addr", ":9898", "Listen address for http server")
	flag.StringVar(&loggerPrefix, "logger-prefix", "", "Logger prefix")
	flag.StringVar(&dbFile, "db-file", "urls.db", "File that is used to persist urls.")
	flag.Parse()
}

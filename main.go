package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilkinulas/youtube-podcast/config"
	"github.com/ilkinulas/youtube-podcast/httpserver"
	"github.com/ilkinulas/youtube-podcast/service"
	"github.com/ilkinulas/youtube-podcast/storage"
	"github.com/ilkinulas/youtube-podcast/version"
)

var (
	configFile string
)

func main() {
	parseFlags()
	logger := log.New(os.Stdout, "", log.LstdFlags)
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		logger.Fatalf("Failed to load config file. %v", err)
	}

	logger.Printf("Starting youtube-podcast app %v", version.GetHumanVersion())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
		<-sigch
		cancel()
	}()

	storage, err := storage.NewSqliteStorage(cfg.DbFile)
	if err != nil {
		logger.Fatalf("Failed to initialize storage, %v", err)
	}

	downloader := service.NewPythonDownloader(logger)
	uploader := service.NewS3Uploader(cfg.S3, logger)

	downloadService := service.NewService(ctx, storage, downloader, uploader, logger)
	go downloadService.Loop()

	server := httpserver.NewServer(ctx, cfg.ListenAddr, logger, storage)
	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Failed to start http server at %v, %v", cfg.ListenAddr, err)
	}
}

func parseFlags() {
	flag.StringVar(&configFile, "config", "config.toml", "path to configuration file.")
	flag.Parse()
}

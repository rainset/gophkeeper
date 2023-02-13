package main

import (
	"context"
	"errors"
	"flag"
	"github.com/caarlos0/env"
	"github.com/rainset/gophkeeper/internal/client/app"
	"github.com/rainset/gophkeeper/internal/client/config"
	"github.com/rainset/gophkeeper/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Config{}
	flag.StringVar(&cfg.ServerAddress, "s", "localhost:8080", "server address")
	flag.StringVar(&cfg.ServerProtocol, "p", "http", "server protocol (http|https)")
	flag.StringVar(&cfg.ClientFolder, "d", "gophkeeper_files", "folder for downloaded files")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		logger.Error(err)

		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	if _, err := os.Stat(cfg.ClientFolder); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(cfg.ClientFolder, os.ModePerm)
		if err != nil {
			logger.Error(err)

			return
		}
	}

	a := app.New(ctx, &cfg)
	a.Run()
}

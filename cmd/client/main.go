package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rainset/gophkeeper/internal/client/app"
	"github.com/rainset/gophkeeper/internal/client/config"
	"github.com/rainset/gophkeeper/pkg/logger"
)

var (
	// BuildVersion is a build version of client application.
	BuildVersion = "N/A"
	// BuildDate is a build date of client application.
	BuildDate = "N/A"
	// BuildCommit is a build commit of client application.
	BuildCommit = "N/A"
)

func main() {

	// print client build version
	fmt.Printf("Build version: %s\n", BuildVersion)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Build commit: %s\n", BuildCommit)

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
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

	a := app.New(ctx, cfg)
	a.Run()
}

package main

import (
	"fmt"
	"github.com/rainset/gophkeeper/internal/client/app"
	"github.com/rainset/gophkeeper/internal/client/config"
	"log"
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

	a := app.New(cfg)
	a.Run()
}

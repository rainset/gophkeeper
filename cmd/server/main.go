package main

import (
	"fmt"
	"github.com/rainset/gophkeeper/internal/server/app"
	"github.com/rainset/gophkeeper/internal/server/config"
)

var (
	// BuildVersion is a build version of server application.
	BuildVersion = "N/A"
	// BuildDate is a build date of server application.
	BuildDate = "N/A"
	// BuildCommit is a build commit of server application.
	BuildCommit = "N/A"
)

func main() {

	// print server build info
	fmt.Printf("Build version: %s\n", BuildVersion)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Build commit: %s\n", BuildCommit)

	cfg, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}
	app.Run(cfg)
}

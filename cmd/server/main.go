package main

import (
	"github.com/rainset/gophkeeper/internal/server/app"
	"github.com/rainset/gophkeeper/internal/server/config"
)

func main() {
	cfg := config.New()
	app.Run(cfg)
}

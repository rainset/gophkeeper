package app

import (
	"context"
	"github.com/rainset/gophkeeper/internal/server/config"
	"github.com/rainset/gophkeeper/internal/server/handler"
	"github.com/rainset/gophkeeper/internal/server/service"
	"github.com/rainset/gophkeeper/internal/server/storage"
	"github.com/rainset/gophkeeper/internal/server/storage/file"
	"log"
	"net/http"
	"testing"
)

func TestNewServer(t *testing.T) {

	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		log.Fatal(err)
	}

	newService := service.New(store, storeFile, cfg)
	newHandler := handler.NewHandler(newService)

	type args struct {
		cfg     *config.Config
		handler http.Handler
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "new server init",
			args: args{
				cfg,
				newHandler.Init(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewServer(tt.args.cfg, tt.args.handler)

			if got.httpServer.Addr != "localhost:8080" {
				t.Errorf("NewServer() = %v, want %v", got, got.httpServer.Addr)
			}
		})
	}
}

func TestRun(t *testing.T) {
	t.Skipped()
}

func TestServer_Run(t *testing.T) {
	t.Skipped()
}

func TestServer_Stop(t *testing.T) {
	t.Skipped()
}

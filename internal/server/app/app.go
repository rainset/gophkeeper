// Package app
package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rainset/gophkeeper/internal/server/config"
	"github.com/rainset/gophkeeper/internal/server/handler"
	"github.com/rainset/gophkeeper/internal/server/service"
	"github.com/rainset/gophkeeper/internal/server/storage"
	"github.com/rainset/gophkeeper/internal/server/storage/file"
	"github.com/rainset/gophkeeper/pkg/logger"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           cfg.ServerAddress,
			Handler:        handler,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServeTLS("./cert/cert.pem", "./cert/private.key")
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// Run Инициализация приложения.
func Run(cfg *config.Config) {
	ctx, cancel := context.WithCancel(context.Background())

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		log.Fatal(err)
	}

	newService := service.New(store, storeFile, cfg)
	newHandler := handler.NewHandler(newService)

	// HTTP Server
	srv := NewServer(cfg, newHandler.Init())

	var wg sync.WaitGroup
	wg.Add(1) // добавляем одну горутину в группу

	// удаление по времени
	go func() {
		defer wg.Done()
		for {
			err := newService.ClearExpiredRefreshTokens(ctx)
			if err != nil {
				logger.Error(err)
				return
			}
			time.Sleep(60 * time.Second)
		}

	}()

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	quit := make(chan os.Signal, 3)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	logger.Info("Shutting down server...")

	wg.Wait()
	store.Close()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	logger.Info("Server exiting")
}

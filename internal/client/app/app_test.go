package app

import (
	"github.com/rainset/gophkeeper/internal/client/config"
	"log"
	"testing"
)

func TestNew(t *testing.T) {

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	t.Run("new app", func(t *testing.T) {
		got := New(cfg)

		if got.cfg != cfg {
			t.Errorf("New() = %v, want %v", got.cfg, cfg)
		}
	})
}

package config

import (
	"github.com/rainset/gophkeeper/internal/server/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_readCommandLineArgs(t *testing.T) {
	type fields struct {
		ServerAddress  string
		ServerProtocol string
		ClientFolder   string
	}
	tests := []struct {
		name   string
		fields fields
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				ServerAddress:  tt.fields.ServerAddress,
				ServerProtocol: tt.fields.ServerProtocol,
				ClientFolder:   tt.fields.ClientFolder,
			}
			c.readCommandLineArgs()
			assert.NotNil(t, c)
		})
	}
}

func TestReadConfig(t *testing.T) {
	t.Run("positive test", func(t *testing.T) {
		got, err := config.ReadConfig()
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
}

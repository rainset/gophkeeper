package config

import (
	"github.com/rainset/gophkeeper/pkg/logger"
	"reflect"
	"testing"
)

func TestReadConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    *Config
		wantErr bool
	}{
		{
			name: "read config with defaults",
			want: &Config{
				ServerAddress:      "localhost:8080",
				DatabaseDsn:        "postgres://root:12345@localhost:5432/gophkeeper",
				FileStorage:        "_file_storage",
				JWTSecretKey:       "secret_key",
				JWTAccessTokenTTL:  "10h",
				JWTRefreshTokenTTL: "720h",
				EnableTLS:          false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadConfig()

			logger.Info(got)
			logger.Info(tt.want)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

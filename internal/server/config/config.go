package config

import (
	"flag"
	"sync"

	"github.com/caarlos0/env/v7"
)

type Config struct {
	ServerAddress      string `env:"SERVER_ADDRESS" envDefault:"localhost:8080" json:"serverAddress"`
	DatabaseDsn        string `env:"DATABASE_DSN" envDefault:"postgres://root:12345@localhost:5432/gophkeeper" json:"databaseDsn"` //nolint:lll
	FileStorage        string `env:"FILE_STORAGE" envDefault:"_file_storage" json:"fileStorage"`
	JWTSecretKey       string `env:"JWT_SECRET_KEY" envDefault:"secret_key" json:"jwtSecretKey"`
	JWTAccessTokenTTL  string `env:"JWT_ACCESS_TOKEN_TTL" envDefault:"2h" json:"jwtAccessTokenTTL"`
	JWTRefreshTokenTTL string `env:"JWT_REFRESH_TOKEN_TTL" envDefault:"720h" json:"jwtRefreshTokenTTL"`
	EnableTLS          bool   `env:"ENABLE_TLS" envDefault:"false" json:"enableTLS"`
}

var once sync.Once //nolint:gochecknoglobals

func (c *Config) readCommandLineArgs() {
	once.Do(func() {
		flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "server and port to listen on")
		flag.StringVar(&c.DatabaseDsn, "d", c.DatabaseDsn, "database dsn")
		flag.StringVar(&c.JWTSecretKey, "j", c.JWTSecretKey, "jwt secret key")
		flag.BoolVar(&c.EnableTLS, "s", c.EnableTLS, "enable secure mode")
		flag.Parse()
	})
}

// ReadConfig merges settings from environment and command line arguments.
func ReadConfig() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)

	if err != nil {
		return nil, err
	}
	cfg.readCommandLineArgs()

	return &cfg, nil
}

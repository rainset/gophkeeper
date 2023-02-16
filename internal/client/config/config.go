package config

import (
	"flag"
	"sync"

	"github.com/caarlos0/env/v7"
)

type Config struct {
	ServerAddress  string `env:"SERVER_ADDRESS" envDefault:"localhost:8080" json:"serverAddress"`
	ServerProtocol string `env:"SERVER_PROTOCOL" envDefault:"http" json:"serverProtocol"`
	ClientFolder   string `env:"CLIENT_FOLDER" envDefault:"gophkeeper_files" json:"clientFolder"`
}

var once sync.Once //nolint:gochecknoglobals

func (c *Config) readCommandLineArgs() {
	once.Do(func() {
		flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "server and port to listen on")
		flag.StringVar(&c.ServerProtocol, "p", c.ServerProtocol, "(http or https) protocol")
		flag.StringVar(&c.ClientFolder, "f", c.ClientFolder, "local client folder")
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

package config

type Config struct {
	ServerAddress  string `env:"SERVER_ADDRESS"`
	ServerProtocol string `env:"SERVER_PROTOCOL"`
	LogFile        string `env:"CLIENT_LOG_FILE"`
	ClientFolder   string `env:"CLIENT_FOLDER"`
}

package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress string
	DatabaseDsn   string
	FileStorage   string

	JWTSecretKey       string
	JWTAccessTokenTTL  time.Duration
	JWTRefreshTokenTTL time.Duration
}

func New() *Config {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
	}

	return &Config{
		ServerAddress:      getEnv("RUN_ADDRESS", ""),
		DatabaseDsn:        getEnv("DATABASE_URI", ""),
		FileStorage:        getEnv("FILE_STORAGE", ""),
		JWTSecretKey:       getEnv("JWT_SECRET_KEY", ""),
		JWTAccessTokenTTL:  getEnvAsTimeDuration("JWT_ACCESS_TOKEN_TTL", "1h"),
		JWTRefreshTokenTTL: getEnvAsTimeDuration("JWT_REFRESH_TOKEN_TTL", "10h"),
	}
}

// Simple helper function to read an environment or return a default value.
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsTimeDuration(name string, defaultVal string) time.Duration {
	valStr := getEnv(name, "")
	if valStr == "" {
		valStr = defaultVal
	}

	d, err := time.ParseDuration(valStr)
	if err != nil {
		return d
	}

	return d
}

// Simple helper function to read an environment variable into integer or return a default value
// func getEnvAsInt(name string, defaultVal int) int {
//	valueStr := getEnv(name, "")
//	if value, err := strconv.Atoi(valueStr); err == nil {
//		return value
//	}
//
//	return defaultVal
//}

// Helper to read an environment variable into a bool or return default value
// func getEnvAsBool(name string, defaultVal bool) bool {
//	valStr := getEnv(name, "")
//	if val, err := strconv.ParseBool(valStr); err == nil {
//		return val
//	}
//
//	return defaultVal
//}

// Helper to read an environment variable into a string slice or return default value
// func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
//	valStr := getEnv(name, "")
//
//	if valStr == "" {
//		return defaultVal
//	}
//
//	val := strings.Split(valStr, sep)
//
//	return val
//}

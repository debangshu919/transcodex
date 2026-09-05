package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Env         string
	DatabaseURL string
	LogLevel    string
}

func MustLoad() Config {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT is required!")
	}

	env := os.Getenv("ENV")
	if env == "" {
		panic("ENV is required!")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		panic("DATABASE_URL is required!")
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	return Config{
		Port:        port,
		Env:         env,
		DatabaseURL: databaseURL,
		LogLevel:    logLevel,
	}
}

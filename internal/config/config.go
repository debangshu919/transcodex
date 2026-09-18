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
	AWSRegion   string
	AWSKey      string
	AWSSecret   string
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

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" && env == "production" {
		panic("AWS Region is required in production!")
	}

	awsKey := os.Getenv("AWS_KEY")
	if awsKey == "" && env == "production" {
		panic("AWS Key is required in production!")
	}

	awsSecret := os.Getenv("AWS_SECRET")
	if awsSecret == "" && env == "production" {
		panic("AWS Key is required in production!")
	}

	if env == "development" {
		awsKey = "test"
		awsSecret = "test"
		awsRegion = "us-east-1"
	}

	return Config{
		Port:        port,
		Env:         env,
		DatabaseURL: databaseURL,
		LogLevel:    logLevel,
		AWSRegion:   awsRegion,
		AWSKey:      awsKey,
		AWSSecret:   awsSecret,
	}
}

package config

import (
	"fmt"
	"os"

	dotEnv "github.com/joho/godotenv"
)

type Config struct {
	GithubToken string
	ServerPort  string
}

func NewConfig() *Config {
	err := dotEnv.Load()
	if err != nil {
		fmt.Printf("error loading .env file: %+v", err.Error())
		panic(err)
	}

	return &Config{
		GithubToken: getEnvOrDefault("GITHUB_TOKEN", ""),
		ServerPort:  getEnvOrDefault("PORT", "8080"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

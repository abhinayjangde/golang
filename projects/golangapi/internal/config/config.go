package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string `env:"DATABASE_URL" envDefault:"mongodb://root:example@localhost:27017/mydb?authSource=admin"`
	Port        string `env:"PORT" envDefault:"8080"`
}

func Load() (Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("No .env file found, using environment variables")
	}

	databaseURL, err := extractEnv("DATABASE_URL")

	if err != nil {
		return Config{}, err
	}

	port, err := extractEnv("PORT")

	if err != nil {
		return Config{}, err
	}

	return Config{
		DatabaseURL: databaseURL,
		Port:        port,
	}, nil
}

func extractEnv(key string) (string, error) {
	val := os.Getenv(key)

	if val == "" {
		return "", fmt.Errorf("Missing env variable")
	}

	return val, nil
}

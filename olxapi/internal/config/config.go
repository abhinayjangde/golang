package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string `env:"PORT" envDefault:"8000"`
	Env         string `env:"ENV" envDefault:"local"`
	DatabaseUrl string `env:"DATABASE_URL"`
}

func MustLoad() Config {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT environment variable is not set")
	}
	env := os.Getenv("ENV")
	if env == "" {
		panic("ENV environment variable is not set")
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		panic("DATABASE_URL environment variable is not set")
	}
	return Config{
		Port:        port,
		Env:         env,
		DatabaseUrl: databaseUrl,
	}
}

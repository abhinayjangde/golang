package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port        string `env:"PORT" envDefault:"8000"`
	Env         string `env:"ENV" envDefault:"local"`
	DatabaseURL string `env:"DATABASE_URL,required"`
	RedisURL    string `env:"REDIS_URL,required"`
	LogFile     string `env:"LOG_FILE" envDefault:"./logs/app.jsonl"`
	JWTSecret   string `env:"JWT_SECRET,required"`
}

func MustLoad() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}

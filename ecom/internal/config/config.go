package config

import (
	"log"
	"time"

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

	AWSRegion    string        `env:"AWS_REGION" envDefault:"ap-south-1"`
	S3Bucket     string        `env:"S3_BUCKET,required"`
	PresignedTTL time.Duration `env:"S3_PRESIGNED_TTL" envDefault:"5m"`
}

func MustLoad() Config {
	godotenv.Load()

	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}

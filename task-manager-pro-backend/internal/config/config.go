package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DatabaseUrl string
	RedisURL    string
	JWTSecret   string
)

func Load() {
	_ = godotenv.Load()

	DatabaseUrl = os.Getenv("DATABASE_URL")
	RedisURL = os.Getenv("REDIS_URL")
	JWTSecret = os.Getenv("JWT_SECRET")

	if DatabaseUrl == "" || RedisURL == "" || JWTSecret == "" {
		log.Fatal("Missing environment variables")
	}
}

package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresURL string
	RedisAddr   string
	RedisPass   string
	RedisDB     int

	RateLimit  int
	RateWindow time.Duration
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Postgres
	pgUser := getEnv("POSTGRES_USER", "postgres")
	pgPass := getEnv("POSTGRES_PASSWORD", "postgres")
	pgDB := getEnv("POSTGRES_DB", "url_shortener")
	pgHost := getEnv("POSTGRES_HOST", "localhost")
	pgPort := getEnv("POSTGRES_PORT", "5432")
	pgURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		pgUser, pgPass, pgHost, pgPort, pgDB)

	// Redis
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisPass := getEnv("REDIS_PASSWORD", "")
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	// Rate limiting
	rateLimit, _ := strconv.Atoi(getEnv("RATE_LIMIT", "5"))
	rateWindowStr := getEnv("RATE_WINDOW", "1m")
	rateWindow, err := time.ParseDuration(rateWindowStr)
	if err != nil {
		log.Fatalf("invalid RATE_WINDOW format: %v", err)
	}

	return &Config{
		PostgresURL: pgURL,
		RedisAddr:   redisAddr,
		RedisPass:   redisPass,
		RedisDB:     redisDB,
		RateLimit:   rateLimit,
		RateWindow:  rateWindow,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

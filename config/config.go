package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort string
	PGDSN    string
}

// вспомогательная функция: взять env с дефолтом
func getEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return def
}

// обязательная env-переменная
func mustEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		log.Fatalf("environment variable %s is required", key)
	}
	return val
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file", err)
	}
	var cfg Config

	// HTTP порт, по умолчанию 8080
	cfg.HTTPPort = getEnv("HTTP_PORT", "8080")
	// Параметры базы
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	user := mustEnv("POSTGRES_USER")
	pass := mustEnv("POSTGRES_PASSWORD")
	db := mustEnv("POSTGRES_DB")

	// Собираем DSN
	cfg.PGDSN = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, pass, host, port, db,
	)

	return &cfg
}

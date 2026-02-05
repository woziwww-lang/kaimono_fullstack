package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type CacheConfig struct {
	TTLSeconds int
}

type AuthConfig struct {
	APIKey string
}

type ServerConfig struct {
	Port         string
	CORSOrigins  []string
	MetricsRoute string
}

type LogConfig struct {
	Level string
}

type Config struct {
	DB     DBConfig
	Redis  RedisConfig
	Cache  CacheConfig
	Auth   AuthConfig
	Server ServerConfig
	Log    LogConfig
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "admin"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "price_comparison"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Cache: CacheConfig{
			TTLSeconds: getEnvInt("CACHE_TTL_SECONDS", 60),
		},
		Auth: AuthConfig{
			APIKey: getEnv("API_KEY", ""),
		},
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			CORSOrigins:  splitCSV(getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:3001")),
			MetricsRoute: getEnv("METRICS_ROUTE", "/metrics"),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Invalid int for %s: %s; using default %d", key, value, defaultValue)
		return defaultValue
	}
	return parsed
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	var cleaned []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return cleaned
}

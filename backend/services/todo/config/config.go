package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort        string
	AppBaseURL string
	AppName   string

	DBHost string
	DBPort string
	DBName string
	DBUser string
	DBPassword string
	DBSSLMode string
	DBMaxOpenConns string
	DBMaxIdleConns string
	DBConnMaxLifetime string

	JWTSecret string
	JWTExpiration string
	JWTRefreshExpiration string

	AllowedOrigins string

	RateLimitRequestPerMinute string
	RateLimitAuthRequestPerMinute string
}

func Load() *Config {
	if os.Getenv("APP_ENV") != "docker" {
		if err := godotenv.Load(".env.local"); err != nil {
			log.Fatalf("error load env: %v", err)
		}
	}

	return &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		AppBaseURL: getEnv("APP_BASE_URL", "http://localhost:8080"),
		AppName: getEnv("APP_NAME", "Donezo"),

		DBHost: getEnv("DB_HOST", ""),
		DBPort: getEnv("DB_PORT", "5432"),
		DBName: getEnv("DB_NAME", "postgres"),
		DBUser: getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBSSLMode: getEnv("DB_SSL_MODE", "require"),
		DBMaxOpenConns: getEnv("DB_MAX_OPEN_CONNS", "25"),
		DBMaxIdleConns: getEnv("DB_MAX_IDLE_CONNS", "10"),
		DBConnMaxLifetime: getEnv("DB_CONN_MAX_LIFETIME", "30m"),
		
		JWTSecret: getEnv("JWT_SECRET", ""),
		JWTExpiration: getEnv("JWT_EXPIRATION", "24h"),
		JWTRefreshExpiration: getEnv("JWT_REFRESH_EXPIRATION", "168h"),

		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:5173"),
		
		RateLimitRequestPerMinute: getEnv("RATE_LIMIT_REQUEST_PER_MINUTE", "60"),
		RateLimitAuthRequestPerMinute: getEnv("RATE_LIMIT_AUTH_REQUEST_PER_MINUTE", "10"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
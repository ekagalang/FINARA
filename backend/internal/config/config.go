package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName       string
	AppEnv        string
	AppPort       string
	DBDriver      string // postgres atau mysql
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	DBCharset     string
	DBParseTime   string
	DBLoc         string
	JWTSecret     string
	JWTExpiration time.Duration
	CORSOrigins   []string
}

func LoadConfig() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		AppName:       getEnv("APP_NAME", "Finara Accounting"),
		AppEnv:        getEnv("APP_ENV", "development"),
		AppPort:       getEnv("APP_PORT", "8080"),
		DBDriver:      getEnv("DB_DRIVER", "postgres"), // Default postgres
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "postgres"),
		DBName:        getEnv("DB_NAME", "finara_db"),
		DBSSLMode:     getEnv("DB_SSLMODE", "disable"),
		DBCharset:     getEnv("DB_CHARSET", "utf8mb4"),
		DBParseTime:   getEnv("DB_PARSE_TIME", "true"),
		DBLoc:         getEnv("DB_LOC", "Local"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiration: 24 * time.Hour,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
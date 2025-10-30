package config

import (
	"os"
)

type Config struct {
	RedisAddr        string
	PostGresUser     string
	PostGresPassword string
	PostGresDB       string
	PostGresHost     string
	PostGresPort     string
	MinioEndpoint    string
	MinioAccessKey   string
	MinioSecretKey   string
	MinioUseSSL      bool
}

func LoadConfig() *Config {
	useSSL := false
	if getEnv("MINIO_USE_SSL", "false") == "true" {
		useSSL = true
	}

	return &Config{
		RedisAddr:        getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		PostGresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostGresPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		PostGresDB:       getEnv("POSTGRES_DB", "rag_chatbot_db"),
		PostGresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostGresPort:     getEnv("POSTGRES_PORT", "5432"),
		MinioEndpoint:    getEnv("MINIO_ENDPOINT", "localhost:9008"),
		MinioAccessKey:   getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey:   getEnv("MINIO_SECRET_KEY", "minioadmin"),
		MinioUseSSL:      useSSL,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

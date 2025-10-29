package config

import (
	"os"
)

type Config struct {
	Port string
	PostGresUser string
	PostGresPassword string
	PostGresDB string
	PostGresHost string
	PostGresPort string
	MinioEndpoint string
	MinioBucket string
	LocalStoragePath string
	RedisAddr string
	JWTSecret string
}

func LoadConfig() *Config {
	return &Config{
		Port:            getEnv("PORT", "8080"),
		PostGresUser:    getEnv("POSTGRES_USER", "postgres"),
		PostGresPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		PostGresDB:      getEnv("POSTGRES_DB", "rag_chatbot_db"),
		PostGresHost:    getEnv("POSTGRES_HOST", "localhost"),
		PostGresPort:    getEnv("POSTGRES_PORT", "5432"),
		MinioEndpoint:   getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioBucket:     getEnv("MINIO_BUCKET", "documents"),
		LocalStoragePath: getEnv("LOCAL_STORAGE_PATH", "./data"),
		RedisAddr:       getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret:       getEnv("JWT_SECRET", "your_jwt_secret_key"),
	}
}


func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
package services

import (
	"backend/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the global GORM DB connection and runs AutoMigrate for models.
func InitDB() error {
	cfg := config.LoadConfig()

	// First, connect to the default 'postgres' database to create our target database if needed
	defaultDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable TimeZone=UTC",
		cfg.PostGresHost,
		cfg.PostGresUser,
		cfg.PostGresPassword,
		cfg.PostGresPort,
	)

	defaultDB, err := gorm.Open(postgres.Open(defaultDSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to default postgres database: %w", err)
	}

	// Check if database exists, create if not
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", cfg.PostGresDB)
	err = defaultDB.Raw(query).Scan(&exists).Error
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		log.Printf("Database '%s' does not exist, creating it...", cfg.PostGresDB)
		createDBSQL := fmt.Sprintf("CREATE DATABASE %s", cfg.PostGresDB)
		err = defaultDB.Exec(createDBSQL).Error
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		log.Printf("Database '%s' created successfully", cfg.PostGresDB)
	} else {
		log.Printf("Database '%s' already exists", cfg.PostGresDB)
	}

	// Close the default database connection
	sqlDB, err := defaultDB.DB()
	if err == nil {
		sqlDB.Close()
	}

	// Now connect to the target database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.PostGresHost,
		cfg.PostGresUser,
		cfg.PostGresPassword,
		cfg.PostGresDB,
		cfg.PostGresPort,
	)

	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to target database: %w", err)
	}

	DB = gdb
	log.Printf("Successfully connected to database '%s'", cfg.PostGresDB)

	return nil
}

// CloseDB closes underlying database connection when possible.
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

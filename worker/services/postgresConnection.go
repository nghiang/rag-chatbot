package services

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(user, password, dbname, host, port string) error {
	// First, connect to the default 'postgres' database to create our target database if needed
	defaultDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable",
		host, user, password, port)

	defaultDB, err := gorm.Open(postgres.Open(defaultDSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to default postgres database: %w", err)
	}

	// Check if database exists, create if not
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", dbname)
	err = defaultDB.Raw(query).Scan(&exists).Error
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		log.Printf("Database '%s' does not exist, creating it...", dbname)
		createDBSQL := fmt.Sprintf("CREATE DATABASE %s", dbname)
		err = defaultDB.Exec(createDBSQL).Error
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		log.Printf("Database '%s' created successfully", dbname)
	} else {
		log.Printf("Database '%s' already exists", dbname)
	}

	// Close the default database connection
	sqlDB, err := defaultDB.DB()
	if err == nil {
		sqlDB.Close()
	}

	// Now connect to the target database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to target database: %w", err)
	}

	log.Printf("Successfully connected to PostgreSQL database '%s'", dbname)
	return nil
}

func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func UpdateDocumentStatus(documentID uint, status string) error {
	result := DB.Table("documents").
		Where("id = ?", documentID).
		Update("embedding_status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update document status: %w", result.Error)
	}

	log.Printf("Updated document ID %d status to: %s", documentID, status)
	return nil
}

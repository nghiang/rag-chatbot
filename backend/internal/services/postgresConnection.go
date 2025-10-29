package services

import (
	"fmt"
	"backend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the global GORM DB connection and runs AutoMigrate for models.
func InitDB() error {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		config.LoadConfig().PostGresHost,
		config.LoadConfig().PostGresUser,
		config.LoadConfig().PostGresPassword,
		config.LoadConfig().PostGresDB,
		config.LoadConfig().PostGresPort,
	)

	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = gdb

	// Run AutoMigrate for models if model package is imported elsewhere.
	// We'll call AutoMigrate from main, after models are defined.
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

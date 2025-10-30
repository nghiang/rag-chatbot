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

// CreateSchemaIfNotExists creates a PostgreSQL schema if it doesn't exist
func CreateSchemaIfNotExists(schemaName string) error {
	// Sanitize schema name to prevent SQL injection
	if !isValidIdentifier(schemaName) {
		return fmt.Errorf("invalid schema name: %s", schemaName)
	}

	sql := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schemaName)
	result := DB.Exec(sql)

	if result.Error != nil {
		return fmt.Errorf("failed to create schema: %w", result.Error)
	}

	log.Printf("Schema '%s' created or already exists", schemaName)
	return nil
}

// CreateTableFromCSV creates a table based on CSV headers
func CreateTableFromCSV(schemaName, tableName string, headers []string) error {
	// Sanitize identifiers
	if !isValidIdentifier(schemaName) {
		return fmt.Errorf("invalid schema name: %s", schemaName)
	}
	if !isValidIdentifier(tableName) {
		return fmt.Errorf("invalid table name: %s", tableName)
	}

	// Build column definitions (all columns as TEXT for simplicity)
	columns := make([]string, 0, len(headers)+1)
	columns = append(columns, "id SERIAL PRIMARY KEY")

	for _, header := range headers {
		// Sanitize column name and use quoted identifier
		sanitizedName := sanitizeColumnName(header)
		columns = append(columns, fmt.Sprintf("\"%s\" TEXT", sanitizedName))
	}

	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s (%s)",
		schemaName, tableName, joinStrings(columns, ", "))

	result := DB.Exec(sql)
	if result.Error != nil {
		return fmt.Errorf("failed to create table: %w", result.Error)
	}

	log.Printf("Table '%s.%s' created successfully with columns: %v", schemaName, tableName, headers)
	return nil
}

// SetTableComment adds a comment/description to a PostgreSQL table
func SetTableComment(schemaName, tableName, comment string) error {
	// Sanitize identifiers
	if !isValidIdentifier(schemaName) {
		return fmt.Errorf("invalid schema name: %s", schemaName)
	}
	if !isValidIdentifier(tableName) {
		return fmt.Errorf("invalid table name: %s", tableName)
	}

	if comment == "" {
		return nil // No comment to set
	}

	// Escape single quotes in comment by doubling them
	escapedComment := ""
	for _, char := range comment {
		if char == '\'' {
			escapedComment += "''"
		} else {
			escapedComment += string(char)
		}
	}

	// PostgreSQL COMMENT ON TABLE command
	// We need to use string literals, not placeholders, for COMMENT statements
	sql := fmt.Sprintf("COMMENT ON TABLE %s.%s IS '%s'", schemaName, tableName, escapedComment)

	result := DB.Exec(sql)
	if result.Error != nil {
		return fmt.Errorf("failed to set table comment: %w", result.Error)
	}

	log.Printf("Comment added to table '%s.%s': %s", schemaName, tableName, comment)
	return nil
}

// InsertCSVRow inserts a row of CSV data into the table
func InsertCSVRow(schemaName, tableName string, headers, values []string) error {
	if len(headers) != len(values) {
		return fmt.Errorf("headers and values count mismatch: %d != %d", len(headers), len(values))
	}

	// Sanitize identifiers
	if !isValidIdentifier(schemaName) {
		return fmt.Errorf("invalid schema name: %s", schemaName)
	}
	if !isValidIdentifier(tableName) {
		return fmt.Errorf("invalid table name: %s", tableName)
	}

	// Build column names with sanitization
	columnNames := make([]string, 0, len(headers))
	for _, header := range headers {
		sanitizedName := sanitizeColumnName(header)
		columnNames = append(columnNames, fmt.Sprintf("\"%s\"", sanitizedName))
	}

	// Build placeholders for values
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	sql := fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES (%s)",
		schemaName, tableName,
		joinStrings(columnNames, ", "),
		joinStrings(placeholders, ", "))

	// Convert values to interface slice
	args := make([]interface{}, len(values))
	for i, v := range values {
		args[i] = v
	}

	result := DB.Exec(sql, args...)
	if result.Error != nil {
		return fmt.Errorf("failed to insert row: %w", result.Error)
	}

	return nil
}

// isValidIdentifier checks if a string is a valid PostgreSQL identifier
func isValidIdentifier(s string) bool {
	if len(s) == 0 || len(s) > 63 {
		return false
	}

	// Must start with letter or underscore
	if !((s[0] >= 'a' && s[0] <= 'z') || (s[0] >= 'A' && s[0] <= 'Z') || s[0] == '_') {
		return false
	}

	// Rest can be letters, digits, or underscores
	for i := 1; i < len(s); i++ {
		c := s[i]
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}

	return true
}

// sanitizeColumnName converts any string to a valid PostgreSQL column name
// This preserves the original characters (including Unicode) since we'll use quoted identifiers
func sanitizeColumnName(s string) string {
	if len(s) == 0 {
		return "col_unnamed"
	}

	// PostgreSQL allows any characters in quoted identifiers except null bytes
	// We just need to escape any double quotes by doubling them
	result := ""
	for _, char := range s {
		if char == '"' {
			result += "\"\""
		} else if char == 0 {
			// Skip null bytes
			continue
		} else {
			result += string(char)
		}
	}

	// Truncate to 63 characters (PostgreSQL identifier limit)
	if len(result) > 63 {
		result = result[:63]
	}

	return result
}

// joinStrings joins strings with a separator (helper to avoid importing strings)
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

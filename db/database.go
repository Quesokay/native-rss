// db/database.go
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB opens the connection and applies the schema
func InitDB(dataSourceName string) error {
	var err error
	
	// The PRAGMA flags here are the secret sauce. 
	// _journal=WAL allows concurrent reads/writes
	// _foreign_keys=on ensures our CASCADE delete works
	dsn := fmt.Sprintf("%s?_journal=WAL&_foreign_keys=on&_timeout=5000", dataSourceName)
	
	DB, err = sql.Open("sqlite3", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Read and execute the schema file
	schema, err := os.ReadFile("db/schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema.sql: %w", err)
	}

	_, err = DB.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	log.Println("Database successfully initialized in WAL mode.")
	return nil
}
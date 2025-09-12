package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// NewPostgresConnection creates a new PostgreSQL connection
func NewPostgresConnection(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Check if connection is alive
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

// createTables creates the necessary tables if they don't exist
func createTables(db *sql.DB) error {
	// Create products table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			price DECIMAL(10, 2) NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Create inventory table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS inventory (
			product_id VARCHAR(255) PRIMARY KEY REFERENCES products(id),
			quantity INT NOT NULL,
			reserved INT NOT NULL DEFAULT 0,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Create reservations table for idempotent reservations
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS reservations (
			order_id VARCHAR(255) NOT NULL,
			product_id VARCHAR(255) NOT NULL REFERENCES products(id),
			quantity INT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			UNIQUE(order_id, product_id)
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

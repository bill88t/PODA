package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDB function
func InitDB(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	if err = createTables(); err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

// createTables function
func createTables() error {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		fname TEXT NOT NULL,
		lname TEXT,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		birthday TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	productsTable := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		brand TEXT NOT NULL,
		price REAL NOT NULL
	);`

	if _, err := db.Exec(usersTable); err != nil {
		return err
	}

	if _, err := db.Exec(productsTable); err != nil {
		return err
	}

	// Seed products if empty
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		seedProducts := `
		INSERT INTO products (name, brand, price) VALUES
			('Urban Runner', 'PODA', 79.99),
			('Classic Loafer', 'PODA', 120.50),
			('Summer Sandal', 'PODA', 45.00);`
		if _, err := db.Exec(seedProducts); err != nil {
			return err
		}
		log.Println("Seeded initial products")
	}

	return nil
}

// CloseDB function
func CloseDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}
}

// GetDB function
func GetDB() *sql.DB {
	return db
}

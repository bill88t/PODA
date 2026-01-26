package main

import (
	"database/sql"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// InitDB opens the sqlite database, configures the underlying pool/pragmas,
// runs AutoMigrate for models, and seeds initial products
func InitDB(dbPath string) error {
	var err error
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// For SQLite it's common to keep pools small to avoid "database is locked"
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := execPragmas(sqlDB); err != nil {
		log.Printf("warning: failed to apply pragmas: %v", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(&UserModel{}, &ProductModel{}); err != nil {
		return err
	}

	// Seed initial products if needed
	if err := seedProducts(); err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

// execPragmas runs SQLite pragmas via a raw Exec on gorm DB
func execPragmas(sqlDB *sql.DB) error {
	gormDB, err := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	// Enable foreign keys
	if err := gormDB.Exec("PRAGMA foreign_keys = ON;").Error; err != nil {
		return err
	}

	if err := gormDB.Exec("PRAGMA journal_mode = WAL;").Error; err != nil {
		log.Printf("warning: couldn't set journal_mode=WAL: %v", err)
	}

	return nil
}

func seedProducts() error {
	var count int64
	if err := db.Model(&ProductModel{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		products := []ProductModel{
			{Name: "Urban Runner", Brand: "PODA", Price: 79.99},
			{Name: "Classic Loafer", Brand: "PODA", Price: 120.50},
			{Name: "Summer Sandal", Brand: "PODA", Price: 45.00},
		}
		if err := db.Create(&products).Error; err != nil {
			return err
		}
		log.Println("Seeded initial products")
	}

	return nil
}

// GetDB returns the gorm DB handle
func GetDB() *gorm.DB {
	return db
}

// CloseDB closes the underlying database/sql DB
func CloseDB() {
	if db == nil {
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("error getting sql.DB: %v", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Printf("error closing database: %v", err)
	}
}

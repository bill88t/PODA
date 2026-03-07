package mariadb

import (
	"fmt"
	"log"
	"os"
	"time"

	"main/drivers"
	model "main/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func dsn() string {
	user := os.Getenv("MARIADB_USER")
	if user == "" {
		user = "poda"
	}
	password := os.Getenv("MARIADB_PASSWORD")
	host := os.Getenv("MARIADB_HOST")
	if host == "" {
		host = "mariadb"
	}
	port := os.Getenv("MARIADB_PORT")
	if port == "" {
		port = "3306"
	}
	database := os.Getenv("MARIADB_DATABASE")
	if database == "" {
		database = "poda"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)
}

func InitDB() error {
	var err error

	for range 10 {
		drivers.Db, err = gorm.Open(mysql.Open(dsn()), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err == nil {
			break
		}
		log.Printf("Database not ready, retrying in 3s: %v", err)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		return fmt.Errorf("failed to connect after retries: %w", err)
	}

	sqlDB, err := drivers.Db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := drivers.Db.AutoMigrate(
		&model.UserModel{},
		&model.AppointmentModel{},
	); err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

func CloseDB() {
	if drivers.Db == nil {
		return
	}

	sqlDB, err := drivers.Db.DB()
	if err != nil {
		log.Printf("error getting sql.DB: %v", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Printf("error closing database: %v", err)
	}
}

package main

import (
	"time"

	"github.com/google/uuid"
)

// UserModel Database type
type UserModel struct {
	ID           uuid.UUID `gorm:"primaryKey"`
	FName        string    `gorm:"not null"`
	LName        string
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Birthday     *string
	CreatedAt    time.Time
}

// TableName UserModel table binding
func (UserModel) TableName() string {
	return "users"
}

// ProductModel Database type
type ProductModel struct {
	ID    uint    `gorm:"primaryKey;autoIncrement"`
	Name  string  `gorm:"not null"`
	Brand string  `gorm:"not null"`
	Price float64 `gorm:"not null"`
}

// TableName ProductModel table binding
func (ProductModel) TableName() string {
	return "products"
}

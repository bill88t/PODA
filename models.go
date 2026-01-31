package main

import (
	"time"

	"github.com/google/uuid"
)

// UserModel Database type
type UserModel struct {
	ID           uuid.UUID `gorm:"primaryKey"`
	Kind         string    `gorm:"not null"`
	FName        string    `gorm:"not null"`
	LName        string    `gorm:"not null"`
	Email        string    `gorm:"uniqueIndex;not null"`
	Phone        string
	PasswordHash string  `gorm:"not null"`
	Birthday     *string `gorm:"not null"`
	Address      string
	Appointments []AppointmentModel `gorm:"foreignKey:UserID"`
}

// TableName UserModel table binding
func (UserModel) TableName() string {
	return "users"
}

// AppointmentModel Database type
type AppointmentModel struct {
	ID       uint      `gorm:"primaryKey;autoIncrement"`
	UserID   uuid.UUID `gorm:"not null;index"`
	Datetime time.Time `gorm:"not null"`
	Kind     string    `gorm:"not null"`
}

// TableName AppointmentModel table binding
func (AppointmentModel) TableName() string {
	return "appointments"
}

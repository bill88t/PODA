package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"main/drivers"
	model "main/models"

	_ "github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User struct for API responses
type User struct {
	ID           uuid.UUID     `json:"id"`
	Kind         string        `json:"kind"`
	Fname        string        `json:"fname"`
	Lname        string        `json:"lname"`
	Email        string        `json:"email"`
	Phone        string        `json:"phone"`
	PasswordHash string        `json:"-"`
	Birthday     time.Time     `json:"birthday"`
	Address      string        `json:"address"`
	Appointments []Appointment `json:"appointments"`
}

// UserResponse struct for JSON output
type UserResponse struct {
	ID           uuid.UUID     `json:"id"`
	Kind         string        `json:"kind"`
	Fname        string        `json:"fname"`
	Lname        string        `json:"lname"`
	Email        string        `json:"email"`
	Phone        string        `json:"phone"`
	Birthday     time.Time     `json:"birthday"`
	Address      string        `json:"address"`
	Appointments []Appointment `json:"appointments"`
}

// hashPassword hashes a password
func hashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

// CreateUser inserts a new user into the database with the given kind ("client" or "barber")
func CreateUser(fname, lname, email, password, kind string, birthday time.Time) (*User, error) {
	id := uuid.New()
	passwordHash := hashPassword(password)

	// Convert birthday to string to match SQLite TEXT
	bdayStr := birthday.Format("2006-01-02")

	userModel := model.UserModel{
		ID:           id,
		Kind:         kind,
		FName:        fname,
		LName:        lname,
		Email:        email,
		PasswordHash: passwordHash,
		Birthday:     &bdayStr,
	}

	if err := drivers.Db.Create(&userModel).Error; err != nil {
		return nil, err
	}

	return &User{
		ID:           id,
		Kind:         kind,
		Fname:        fname,
		Lname:        lname,
		Email:        email,
		Birthday:     birthday,
		Appointments: []Appointment{},
	}, nil
}

// GetUserByEmail fetches a user by email
func GetUserByEmail(email string) (*User, error) {
	var user model.UserModel
	if err := drivers.Db.Preload("Appointments").Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}

	// Convert birthday string to time.Time
	birthday := time.Time{}
	if user.Birthday != nil {
		parsed, err := time.Parse("2006-01-02", *user.Birthday)
		if err == nil {
			birthday = parsed
		}
	}

	appointments := make([]Appointment, len(user.Appointments))
	for i, appmnt := range user.Appointments {
		appointments[i] = Appointment{
			ID:       appmnt.ID,
			Datetime: appmnt.Datetime,
			Kind:     appmnt.Kind,
		}
	}

	return &User{
		ID:           user.ID,
		Kind:         user.Kind,
		Fname:        user.FName,
		Lname:        user.LName,
		Email:        user.Email,
		Phone:        user.Phone,
		PasswordHash: user.PasswordHash,
		Birthday:     birthday,
		Address:      user.Address,
		Appointments: appointments,
	}, nil
}

// GetUserByID fetches a user by ID
func GetUserByID(id string) (*User, error) {
	var user model.UserModel
	if err := drivers.Db.Preload("Appointments").First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}

	birthday := time.Time{}
	if user.Birthday != nil {
		parsed, err := time.Parse("2006-01-02", *user.Birthday)
		if err == nil {
			birthday = parsed
		}
	}

	appointments := make([]Appointment, len(user.Appointments))
	for i, appmnt := range user.Appointments {
		appointments[i] = Appointment{
			ID:       appmnt.ID,
			Datetime: appmnt.Datetime,
			Kind:     appmnt.Kind,
		}
	}

	return &User{
		ID:           user.ID,
		Kind:         user.Kind,
		Fname:        user.FName,
		Lname:        user.LName,
		Email:        user.Email,
		Phone:        user.Phone,
		PasswordHash: user.PasswordHash,
		Birthday:     birthday,
		Address:      user.Address,
		Appointments: appointments,
	}, nil
}

// UpdateUserPassword updates the password hash for a user
func UpdateUserPassword(id string, password string) error {
	return drivers.Db.Model(&model.UserModel{}).
		Where("id = ?", id).
		Update("password_hash", hashPassword(password)).Error
}

// UpdateUserInfo updates the name and birthday for a user
func UpdateUserInfo(id string, fname, lname string, birthday time.Time) error {
	bdayStr := birthday.Format("2006-01-02")
	return drivers.Db.Model(&model.UserModel{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"f_name":   fname,
			"l_name":   lname,
			"birthday": bdayStr,
		}).Error
}

// UpdateUserContact updates the email and phone for a user
func UpdateUserContact(id string, email, phone string) error {
	return drivers.Db.Model(&model.UserModel{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"email": email,
			"phone": phone,
		}).Error
}

// GetAllAppointments fetches every appointment across all users
func GetAllAppointments() ([]model.AppointmentModel, error) {
	var appointments []model.AppointmentModel
	if err := drivers.Db.Find(&appointments).Error; err != nil {
		return nil, err
	}
	return appointments, nil
}

// EmailExists checks if an email is already registered
func EmailExists(email string) bool {
	var count int64
	if err := drivers.Db.Model(
		&model.UserModel{},
	).Where("email = ?", email).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}

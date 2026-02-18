package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"main/drivers"
	model "main/models"

	"github.com/gofiber/fiber/v2"
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

// CreateUser inserts a new user into the database
func CreateUser(fname, lname, email, password string, birthday time.Time) (*User, error) {
	id := uuid.New()
	passwordHash := hashPassword(password)

	// Convert birthday to string to match SQLite TEXT
	bdayStr := birthday.Format("2006-01-02")

	userModel := model.UserModel{
		ID:           id,
		Kind:         "client",
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
		Kind:         "client",
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

// SignUp handler
func SignUp(ctx *fiber.Ctx) error {
	req := new(struct {
		Fname    string `json:"fname"`
		Lname    string `json:"lname"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Birthday string `json:"birthday"`
	})
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	if req.Email == "" || req.Password == "" || req.Fname == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "First name, email, and password are required"})
	}

	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid birthday format"})
	}

	if EmailExists(req.Email) {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already exists"})
	}

	user, err := CreateUser(req.Fname, req.Lname, req.Email, req.Password, birthday)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	token, err := GenerateToken(user.ID, user.Email)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user": UserResponse{
			ID:           user.ID,
			Kind:         user.Kind,
			Fname:        user.Fname,
			Lname:        user.Lname,
			Email:        user.Email,
			Phone:        user.Phone,
			Birthday:     user.Birthday,
			Address:      user.Address,
			Appointments: []Appointment{},
		},
		"token": token,
	})
}

// Login handler
func Login(ctx *fiber.Ctx) error {
	req := new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	user, err := GetUserByEmail(req.Email)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	if user.PasswordHash != hashPassword(req.Password) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	token, err := GenerateToken(user.ID, user.Email)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return ctx.JSON(fiber.Map{
		"user": UserResponse{
			ID:           user.ID,
			Kind:         user.Kind,
			Fname:        user.Fname,
			Lname:        user.Lname,
			Email:        user.Email,
			Phone:        user.Phone,
			Birthday:     user.Birthday,
			Address:      user.Address,
			Appointments: user.Appointments,
		},
		"token": token,
	})
}

// GetProfile handler
func GetProfile(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uuid.UUID)

	user, err := GetUserByID(userID.String())
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return ctx.JSON(UserResponse{
		ID:           user.ID,
		Kind:         user.Kind,
		Fname:        user.Fname,
		Lname:        user.Lname,
		Email:        user.Email,
		Phone:        user.Phone,
		Birthday:     user.Birthday,
		Address:      user.Address,
		Appointments: user.Appointments,
	})
}

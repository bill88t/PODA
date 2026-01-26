package main

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User struct for API responses
type User struct {
	ID           string    `json:"id"`
	Fname        string    `json:"fname"`
	Lname        string    `json:"lname"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Birthday     time.Time `json:"birthday"`
	CreatedAt    time.Time `json:"created_at"`
}

// UserResponse struct for JSON output
type UserResponse struct {
	ID       string    `json:"id"`
	Fname    string    `json:"fname"`
	Lname    string    `json:"lname"`
	Email    string    `json:"email"`
	Birthday time.Time `json:"birthday"`
}

// hashPassword hashes a password
func hashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

// CreateUser inserts a new user into the database
func CreateUser(fname, lname, email, password string, birthday time.Time) (*User, error) {
	id := uuid.New().String()
	passwordHash := hashPassword(password)

	// Convert birthday to string to match SQLite TEXT
	bdayStr := birthday.Format("2006-01-02")

	userModel := UserModel{
		ID:           id,
		FName:        fname,
		LName:        lname,
		Email:        email,
		PasswordHash: passwordHash,
		Birthday:     &bdayStr,
	}

	if err := db.Create(&userModel).Error; err != nil {
		return nil, err
	}

	return &User{
		ID:       id,
		Fname:    fname,
		Lname:    lname,
		Email:    email,
		Birthday: birthday,
	}, nil
}

// GetUserByEmail fetches a user by email
func GetUserByEmail(email string) (*User, error) {
	var u UserModel
	if err := db.Where("email = ?", email).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}

	// Convert birthday string to time.Time
	birthday := time.Time{}
	if u.Birthday != nil {
		parsed, err := time.Parse("2006-01-02", *u.Birthday)
		if err == nil {
			birthday = parsed
		}
	}

	return &User{
		ID:           u.ID,
		Fname:        u.FName,
		Lname:        u.LName,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Birthday:     birthday,
		CreatedAt:    u.CreatedAt,
	}, nil
}

// GetUserByID fetches a user by ID
func GetUserByID(id string) (*User, error) {
	var u UserModel
	if err := db.First(&u, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}

	birthday := time.Time{}
	if u.Birthday != nil {
		parsed, err := time.Parse("2006-01-02", *u.Birthday)
		if err == nil {
			birthday = parsed
		}
	}

	return &User{
		ID:           u.ID,
		Fname:        u.FName,
		Lname:        u.LName,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Birthday:     birthday,
		CreatedAt:    u.CreatedAt,
	}, nil
}

// EmailExists checks if an email is already registered
func EmailExists(email string) bool {
	var count int64
	if err := db.Model(&UserModel{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}

// SignUp handler
func SignUp(c *fiber.Ctx) error {
	req := new(struct {
		Fname    string `json:"fname"`
		Lname    string `json:"lname"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Birthday string `json:"birthday"`
	})
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	if req.Email == "" || req.Password == "" || req.Fname == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "First name, email, and password are required"})
	}

	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid birthday format"})
	}

	if EmailExists(req.Email) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already exists"})
	}

	user, err := CreateUser(req.Fname, req.Lname, req.Email, req.Password, birthday)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	token, err := GenerateToken(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user": UserResponse{
			ID:       user.ID,
			Fname:    user.Fname,
			Lname:    user.Lname,
			Email:    user.Email,
			Birthday: user.Birthday,
		},
		"token": token,
	})
}

// Login handler
func Login(c *fiber.Ctx) error {
	req := new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	user, err := GetUserByEmail(req.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	if user.PasswordHash != hashPassword(req.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	token, err := GenerateToken(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"user": UserResponse{
			ID:       user.ID,
			Fname:    user.Fname,
			Lname:    user.Lname,
			Email:    user.Email,
			Birthday: user.Birthday,
		},
		"token": token,
	})
}

// GetProfile handler
func GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	user, err := GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(UserResponse{
		ID:       user.ID,
		Fname:    user.Fname,
		Lname:    user.Lname,
		Email:    user.Email,
		Birthday: user.Birthday,
	})
}

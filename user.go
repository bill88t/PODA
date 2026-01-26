package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// User struct
type User struct {
	ID           string    `json:"id"`
	Fname        string    `json:"fname"`
	Lname        string    `json:"lname"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Birthday     time.Time `json:"birthday"`
	CreatedAt    time.Time `json:"created_at"`
}

// UserResponse struct
type UserResponse struct {
	ID       string    `json:"id"`
	Fname    string    `json:"fname"`
	Lname    string    `json:"lname"`
	Email    string    `json:"email"`
	Birthday time.Time `json:"birthday"`
}

// hashPassword function
func hashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

// CreateUser function
func CreateUser(fname, lname, email, password string, birthday time.Time) (*User, error) {
	id := uuid.New().String()
	passwordHash := hashPassword(password)

	_, err := db.Exec(
		"INSERT INTO users (id, fname, lname, email, password_hash, birthday) VALUES (?, ?, ?, ?, ?, ?)",
		id, fname, lname, email, passwordHash, birthday.Format("2006-01-02"),
	)
	if err != nil {
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

// GetUserByEmail function
func GetUserByEmail(email string) (*User, error) {
	var user User
	var birthdayStr sql.NullString
	var createdAt sql.NullString

	err := db.QueryRow(
		"SELECT id, fname, lname, email, password_hash, birthday, created_at FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Fname, &user.Lname, &user.Email, &user.PasswordHash, &birthdayStr, &createdAt)
	if err != nil {
		return nil, err
	}

	if birthdayStr.Valid {
		user.Birthday, _ = time.Parse("2006-01-02", birthdayStr.String)
	}

	return &user, nil
}

// GetUserByID function
func GetUserByID(id string) (*User, error) {
	var user User
	var birthdayStr sql.NullString
	var createdAt sql.NullString

	err := db.QueryRow(
		"SELECT id, fname, lname, email, password_hash, birthday, created_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Fname, &user.Lname, &user.Email, &user.PasswordHash, &birthdayStr, &createdAt)
	if err != nil {
		return nil, err
	}

	if birthdayStr.Valid {
		user.Birthday, _ = time.Parse("2006-01-02", birthdayStr.String)
	}

	return &user, nil
}

// EmailExists function
func EmailExists(email string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	return err == nil && count > 0
}

// SignUp function
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

// Login function
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

// GetProfile function
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

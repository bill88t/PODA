package handlers

import (
	"time"

	"main/middleware"

	"github.com/gofiber/fiber/v2"
)

// Admintoken static token for barber registration
const Admintoken = "poda-barber-admin-token"

// SignUp main registration function
func SignUp(ctx *fiber.Ctx) error {
	req := new(struct {
		Fname      string `json:"fname"`
		Lname      string `json:"lname"`
		Email      string `json:"email"`
		Password   string `json:"password"`
		Birthday   string `json:"birthday"`
		Admintoken string `json:"admin_token"`
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

	kind := "client"
	if req.Admintoken == Admintoken {
		kind = "barber"
	}

	user, err := CreateUser(req.Fname, req.Lname, req.Email, req.Password, kind, birthday)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	token, err := middleware.GenerateToken(user.ID, user.Email)
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

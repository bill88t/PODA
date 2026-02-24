package handlers

import (
	"main/middleware"

	"github.com/gofiber/fiber/v2"
)

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

	token, err := middleware.GenerateToken(user.ID, user.Email)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Error"})
	}

	ctx.Response().Header.Add("Authorization", "Bearer "+token)

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

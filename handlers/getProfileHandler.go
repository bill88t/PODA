package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

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

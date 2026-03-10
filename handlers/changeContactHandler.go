package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ChangeContact updates the email and phone for the authenticated user
func ChangeContact(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uuid.UUID)

	req := new(struct {
		Email string `json:"email"`
		Phone string `json:"phone"`
	})
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	if req.Email == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
	}

	if err := UpdateUserContact(userID.String(), req.Email, req.Phone); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update contact"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

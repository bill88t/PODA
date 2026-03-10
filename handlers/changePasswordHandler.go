package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ChangePassword updates the password for the authenticated user
func ChangePassword(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uuid.UUID)

	req := new(struct {
		Password string `json:"password"`
	})
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	if req.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password is required"})
	}

	if err := UpdateUserPassword(userID.String(), req.Password); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update password"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

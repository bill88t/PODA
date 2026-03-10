package middleware

import (
	"main/drivers"

	model "main/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// BarberMiddleware ensures the authenticated user has kind="barber"
func BarberMiddleware(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("userID").(uuid.UUID)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var user model.UserModel
	if err := drivers.Db.First(&user, "id = ?", userID).Error; err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if user.Kind != "barber" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
	}

	return ctx.Next()
}

package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ChangeInfo updates the name and birthday for the authenticated user
func ChangeInfo(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uuid.UUID)

	req := new(struct {
		Fname    string `json:"fname"`
		Lname    string `json:"lname"`
		Birthday string `json:"birthday"`
	})
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	if req.Fname == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "First name is required"})
	}

	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid birthday format"})
	}

	if err := UpdateUserInfo(userID.String(), req.Fname, req.Lname, birthday); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update info"})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

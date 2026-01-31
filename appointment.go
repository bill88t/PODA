package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Appointment struct for API representation
type Appointment struct {
	ID       uint      `json:"id"`
	Datetime time.Time `json:"datetime"`
	Kind     string    `json:"kind"`
}

// GetUserAppointments fetches all appointments for the authenticated user
func GetUserAppointments(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	userUUID, err := uuid.Parse(userID.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var appointments []AppointmentModel
	if err := db.Where("user_id = ?", userUUID).Find(&appointments).Error; err != nil {
		log.Printf("failed to fetch appointments: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch appointments"})
	}

	result := make([]Appointment, len(appointments))
	for i, a := range appointments {
		result[i] = Appointment{
			ID:       a.ID,
			Datetime: a.Datetime,
			Kind:     a.Kind,
		}
	}

	return c.JSON(result)
}

// GetUserAppointmentByID fetches a single appointment by ID for the authenticated user
func GetUserAppointmentByID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	userUUID, err := uuid.Parse(userID.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	appointmentID := c.Params("id")

	var appointment AppointmentModel
	if err := db.Where("id = ? AND user_id = ?", appointmentID, userUUID).First(&appointment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Appointment not found"})
		}
		log.Printf("failed to fetch appointment by id: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch appointment"})
	}

	return c.JSON(Appointment{
		ID:       appointment.ID,
		Datetime: appointment.Datetime,
		Kind:     appointment.Kind,
	})
}

// CreateAppointment creates a new appointment for the authenticated user
func CreateAppointment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	userUUID, err := uuid.Parse(userID.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	req := new(struct {
		Datetime string `json:"datetime"`
		Kind     string `json:"kind"`
	})
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	if req.Datetime == "" || req.Kind == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datetime and kind are required"})
	}

	datetime, err := time.Parse(time.RFC3339, req.Datetime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid datetime format"})
	}

	appointment := AppointmentModel{
		UserID:   userUUID,
		Datetime: datetime,
		Kind:     req.Kind,
	}

	if err := db.Create(&appointment).Error; err != nil {
		log.Printf("failed to create appointment: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create appointment"})
	}

	return c.Status(fiber.StatusCreated).JSON(Appointment{
		ID:       appointment.ID,
		Datetime: appointment.Datetime,
		Kind:     appointment.Kind,
	})
}

// UpdateAppointment updates an existing appointment for the authenticated user
func UpdateAppointment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	userUUID, err := uuid.Parse(userID.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	appointmentID := c.Params("id")

	var appointment AppointmentModel
	if err := db.Where("id = ? AND user_id = ?", appointmentID, userUUID).First(&appointment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Appointment not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch appointment"})
	}

	req := new(struct {
		Datetime string `json:"datetime"`
		Kind     string `json:"kind"`
	})
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	if req.Datetime != "" {
		datetime, err := time.Parse(time.RFC3339, req.Datetime)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid datetime format"})
		}
		appointment.Datetime = datetime
	}

	if req.Kind != "" {
		appointment.Kind = req.Kind
	}

	if err := db.Save(&appointment).Error; err != nil {
		log.Printf("failed to update appointment: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update appointment"})
	}

	return c.JSON(Appointment{
		ID:       appointment.ID,
		Datetime: appointment.Datetime,
		Kind:     appointment.Kind,
	})
}

// DeleteAppointment deletes an appointment for the authenticated user
func DeleteAppointment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	userUUID, err := uuid.Parse(userID.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	appointmentID := c.Params("id")

	var appointment AppointmentModel
	if err := db.Where("id = ? AND user_id = ?", appointmentID, userUUID).First(&appointment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Appointment not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch appointment"})
	}

	if err := db.Delete(&appointment).Error; err != nil {
		log.Printf("failed to delete appointment: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete appointment"})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

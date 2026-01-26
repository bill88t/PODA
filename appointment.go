package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Appointment struct stays as the API representation
type Appointment struct {
	ID       uuid.UUID `json:"id"`
	Datetime time.Time `json:"datetime"`
	Kind     string    `json:"kind"`
}

// GetAllAppointments fetches all products from the database using GORM
func GetAllAppointments(c *fiber.Ctx) error {
	var dbAppointments []AppointmentModel
	if err := db.Find(&dbAppointments).Error; err != nil {
		log.Printf("failed to fetch appointments: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch appointments"})
	}

	appointments := make([]Appointment, len(dbAppointments))
	for i, p := range dbAppointments {
		appointments[i] = Appointment(p)
	}

	return c.JSON(appointments)
}

// GetAppointmentByID fetches a single product by ID using GORM
func GetAppointmentByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var dbAppointments AppointmentModel
	if err := db.First(&dbAppointments, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
		}
		log.Printf("failed to fetch product by id: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch product"})
	}

	appointment := AppointmentModel{
		ID:       dbAppointments.ID,
		Datetime: dbAppointments.Datetime,
		Kind:     dbAppointments.Kind,
	}

	return c.JSON(appointment)
}

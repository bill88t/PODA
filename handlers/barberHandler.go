package handlers

import (
	"log"

	"main/drivers"
	model "main/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// BarberAppointment is an appointment with its owner's user ID
type BarberAppointment struct {
	Appointment
	UserID string `json:"user_id"`
}

// GetAllAppointmentsHandler returns every appointment across all users (barber only)
func GetAllAppointmentsHandler(ctx *fiber.Ctx) error {
	appointments, err := GetAllAppointments()
	if err != nil {
		log.Printf("failed to fetch all appointments: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch appointments"})
	}

	result := make([]BarberAppointment, len(appointments))
	for i, a := range appointments {
		result[i] = BarberAppointment{
			Appointment: Appointment{
				ID:       a.ID,
				Datetime: a.Datetime,
				Kind:     a.Kind,
			},
			UserID: a.UserID.String(),
		}
	}

	return ctx.JSON(result)
}

// CancelAnyAppointment deletes any appointment by ID (barber only)
func CancelAnyAppointment(ctx *fiber.Ctx) error {
	appointmentID := ctx.Params("id")

	var appointment model.AppointmentModel
	if err := drivers.Db.First(&appointment, "id = ?", appointmentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Appointment not found"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch appointment"})
	}

	if err := drivers.Db.Delete(&appointment).Error; err != nil {
		log.Printf("failed to delete appointment: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to cancel appointment"})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

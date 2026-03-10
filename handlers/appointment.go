package handlers

import (
	"log"
	"time"

	"main/drivers"
	model "main/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Appointment struct for API representation
type Appointment struct {
	ID       uint      `json:"id"`
	UserID   string    `json:"user_id,omitempty"`
	Datetime time.Time `json:"datetime"`
	Kind     string    `json:"kind"`
}

// GetUserAppointments fetches appointments for the authenticated user.
// Barbers receive all appointments across all users.
func GetUserAppointments(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uuid.UUID)

	var appointments []model.AppointmentModel

	if IsBarber(userID) {
		if err := drivers.Db.Find(&appointments).Error; err != nil {
			log.Printf("failed to fetch all appointments: %v", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch appointments"})
		}
	} else {
		if err := drivers.Db.Where("user_id = ?", userID).Find(&appointments).Error; err != nil {
			log.Printf("failed to fetch appointments: %v", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch appointments"})
		}
	}

	isBarber := IsBarber(userID)
	result := make([]Appointment, len(appointments))
	for i, a := range appointments {
		appt := Appointment{
			ID:       a.ID,
			Datetime: a.Datetime,
			Kind:     a.Kind,
		}
		if isBarber {
			appt.UserID = a.UserID.String()
		}
		result[i] = appt
	}

	return ctx.JSON(result)
}

// GetUserAppointmentByID fetches a single appointment by ID.
// Barbers can fetch any appointment; clients are restricted to their own.
func GetUserAppointmentByID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	appointmentID := c.Params("id")

	var appointment model.AppointmentModel
	var err error

	if IsBarber(userID) {
		err = drivers.Db.First(&appointment, "id = ?", appointmentID).Error
	} else {
		err = drivers.Db.Where("id = ? AND user_id = ?", appointmentID, userID).First(&appointment).Error
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Appointment not found"})
		}
		log.Printf("failed to fetch appointment by id: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch appointment"})
	}

	appt := Appointment{
		ID:       appointment.ID,
		Datetime: appointment.Datetime,
		Kind:     appointment.Kind,
	}
	if IsBarber(userID) {
		appt.UserID = appointment.UserID.String()
	}

	return c.JSON(appt)
}

// CreateAppointment creates a new appointment for the authenticated user
func CreateAppointment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

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

	appointment := model.AppointmentModel{
		UserID:   userID,
		Datetime: datetime,
		Kind:     req.Kind,
	}

	if err := drivers.Db.Create(&appointment).Error; err != nil {
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

	appointmentID := c.Params("id")

	var appointment model.AppointmentModel
	if err := drivers.Db.Where("id = ? AND user_id = ?", appointmentID, userID).First(&appointment).Error; err != nil {
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

	if err := drivers.Db.Save(&appointment).Error; err != nil {
		log.Printf("failed to update appointment: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update appointment"})
	}

	return c.JSON(Appointment{
		ID:       appointment.ID,
		Datetime: appointment.Datetime,
		Kind:     appointment.Kind,
	})
}

// DeleteAppointment deletes an appointment.
// Barbers can delete any appointment; clients are restricted to their own.
func DeleteAppointment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	appointmentID := c.Params("id")

	var appointment model.AppointmentModel
	var err error

	if IsBarber(userID) {
		err = drivers.Db.First(&appointment, "id = ?", appointmentID).Error
	} else {
		err = drivers.Db.Where("id = ? AND user_id = ?", appointmentID, userID).First(&appointment).Error
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Appointment not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch appointment"})
	}

	if err := drivers.Db.Delete(&appointment).Error; err != nil {
		log.Printf("failed to delete appointment: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete appointment"})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

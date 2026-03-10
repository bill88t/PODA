package handlers

import (
	"main/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// APIHandler backend function
func APIHandler(app *fiber.App, port string) *fiber.App {
	app.Use(middleware.CorsMiddleware(port))
	app.Use(logger.New())

	api := app.Group("/api")
	api.Use(middleware.LimiterMiddleware())

	v1 := api.Group("/v1")

	// Public routes
	v1.Post("/users/signup", SignUp)
	v1.Post("/users/login", Login)

	// User self-service routes (auth required)
	users := v1.Group("/users")
	users.Use(middleware.AuthMiddleware)
	users.Post("/changepassword", ChangePassword)
	users.Post("/changeinfo", ChangeInfo)
	users.Post("/changecontact", ChangeContact)

	// Protected user profile routes
	protected := v1.Group("/profile")
	protected.Use(middleware.AuthMiddleware)

	protected.Get("/", GetProfile)

	// Appointments for the current user
	appointments := protected.Group("/appointments")
	appointments.Get("/", GetUserAppointments)
	appointments.Get("/:id", GetUserAppointmentByID)
	appointments.Post("/", CreateAppointment)
	appointments.Put("/:id", UpdateAppointment)
	appointments.Delete("/:id", DeleteAppointment)

	return app
}

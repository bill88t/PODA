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

	// api
	api := app.Group("/api")
	api.Use(middleware.LimiterMiddleware())

	// api/v1
	v1 := api.Group("/v1")

	// Public routes
	// api/v1/users
	users := v1.Group("/users")

	users.Post("/signup", SignUp)
	users.Post("/login", Login)

	// User self-service routes (auth required)
	users.Use(middleware.AuthMiddleware)
	users.Post("/changepassword", ChangePassword)
	users.Post("/changeinfo", ChangeInfo)
	users.Post("/changecontact", ChangeContact)

	// Protected user profile routes
	// api/v1/profile
	profile := v1.Group("/profile")
	profile.Use(middleware.AuthMiddleware)

	profile.Get("/", GetProfile)

	// Appointments for the current user
	// api/v1/profile/appointment
	appointments := profile.Group("/appointments")
	appointments.Get("/", GetUserAppointments)
	appointments.Get("/:id", GetUserAppointmentByID)
	appointments.Post("/", CreateAppointment)
	appointments.Put("/:id", UpdateAppointment)
	appointments.Delete("/:id", DeleteAppointment)

	return app
}

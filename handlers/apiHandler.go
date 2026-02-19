package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func ApiHandler(app *fiber.App) *fiber.App {
	app.Use(limiter.New(limiter.Config{
		Max:               15,
		Expiration:        15 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://0.0.0.0:5174",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(logger.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Public routes
	v1.Post("/users/signup", SignUp)
	v1.Post("/users/login", Login)

	// Protected routes
	protected := v1.Group("/profile")
	protected.Use(AuthMiddleware)

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

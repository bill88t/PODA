package main

import (
	"log"
	"os"
	"time"

	"main/drivers"
	"main/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	if err := drivers.InitDB("./poda.db"); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer drivers.CloseDB()

	app := fiber.New()

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
	v1.Post("/users/signup", handlers.SignUp)
	v1.Post("/users/login", handlers.Login)

	// Protected routes
	protected := v1.Group("/profile")
	protected.Use(handlers.AuthMiddleware)

	protected.Get("/", handlers.GetProfile)

	// Appointments for the current user
	appointments := protected.Group("/appointments")
	appointments.Get("/", handlers.GetUserAppointments)
	appointments.Get("/:id", handlers.GetUserAppointmentByID)
	appointments.Post("/", handlers.CreateAppointment)
	appointments.Put("/:id", handlers.UpdateAppointment)
	appointments.Delete("/:id", handlers.DeleteAppointment)

	app.Static("/", "./frontend/dist", fiber.Static{
		Compress: true,
	})
	app.Get("*", func(c *fiber.Ctx) error { return c.SendFile("./frontend/dist/index.html") })

	port := os.Getenv("PORT")
	if port == "" {
		port = "5173"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

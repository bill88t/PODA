package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	if err := InitDB("./poda.db"); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer CloseDB()

	app := fiber.New()

	app.Use(limiter.New(limiter.Config{
		Max:               15,
		Expiration:        15 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://0.0.0.0:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(logger.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Public routes
	v1.Get("/products", GetAllProducts)
	v1.Get("/products/:id", GetProductByID)
	v1.Post("/users/signup", SignUp)
	v1.Post("/users/login", Login)

	// Protected routes
	protected := v1.Group("/protected")
	protected.Use(AuthMiddleware)
	protected.Get("/profile", GetProfile)

	app.Static("/", "./frontend/dist")
	app.Get("*", func(c *fiber.Ctx) error { return c.SendFile("./frontend/dist/index.html") })

	port := os.Getenv("PORT")
	if port == "" {
		port = "5173"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

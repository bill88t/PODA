package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func main() {
	app := fiber.New()
	app.Use(limiter.New(limiter.Config{
		Max:               5,
		Expiration:        15 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	// Folder to set up the assets
	app.Static("/", "./frontend/dist")
	// Fallback for React if a file failed
	app.Static("*", "./frontend/dist/index.html")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal(app.Listen(":" + "5173"))
	} else {
		log.Fatal(app.Listen(":" + port))
	}
}

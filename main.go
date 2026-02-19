package main

import (
	"log"
	"os"

	_ "main/drivers"
	drivers "main/drivers/sqlite"
	"main/handlers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	if err := drivers.InitDB("./poda.db"); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer drivers.CloseDB()

	handlers.ApiHandler(app)

	handlers.StaticHandler(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5173"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Creates a new Fiber instance.
	app := fiber.New()

	// https://github.com/gofiber/recipes/blob/master/react-router/cmd/react-router/main.go
	// Folder to set up the assets
	app.Static("/", "./frontend/dist")
	// Fallvack for React if a file failed
	app.Static("*", "./frontend/dist/index.html")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal(app.Listen(":" + "5173"))
	} else {
		log.Fatal(app.Listen(":" + port))
	}
}

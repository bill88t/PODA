package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Product struct
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Brand string  `json:"brand"`
	Price float64 `json:"price"`
}

var productsDB = []Product{
	{ID: 1, Name: "Urban Runner", Brand: "PODA", Price: 79.99},
	{ID: 2, Name: "Classic Loafer", Brand: "PODA", Price: 120.50},
	{ID: 3, Name: "Summer Sandal", Brand: "PODA", Price: 45.00},
}

func getAllProducts(c *fiber.Ctx) error {
	return c.JSON(productsDB)
}

func getProductByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	for _, p := range productsDB {
		if p.ID == id {
			return c.JSON(p)
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "Product not found",
	})
}

func main() {
	app := fiber.New()

	app.Use(limiter.New(limiter.Config{
		Max:               15,
		Expiration:        15 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://0.0.0.0:5173",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Use(logger.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/products", getAllProducts)
	v1.Get("/products/:id", getProductByID)

	app.Static("/", "./frontend/dist")

	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile("./frontend/dist/index.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5173"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

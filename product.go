package main

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Product struct
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Brand string  `json:"brand"`
	Price float64 `json:"price"`
}

// GetAllProducts function
func GetAllProducts(c *fiber.Ctx) error {
	rows, err := db.Query("SELECT id, name, brand, price FROM products")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch products"})
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Brand, &p.Price); err != nil {
			continue
		}
		products = append(products, p)
	}

	if products == nil {
		products = []Product{}
	}

	return c.JSON(products)
}

// GetProductByID function
func GetProductByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var p Product
	err = db.QueryRow("SELECT id, name, brand, price FROM products WHERE id = ? ", id).Scan(&p.ID, &p.Name, &p.Brand, &p.Price)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	return c.JSON(p)
}

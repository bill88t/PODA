package main

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Product struct stays as the API representation
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Brand string  `json:"brand"`
	Price float64 `json:"price"`
}

// GetAllProducts fetches all products from the database using GORM
func GetAllProducts(c *fiber.Ctx) error {
	var dbProducts []ProductModel
	if err := db.Find(&dbProducts).Error; err != nil {
		log.Printf("failed to fetch products: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch products"})
	}

	products := make([]Product, len(dbProducts))
	for i, p := range dbProducts {
		products[i] = Product{
			ID:    int(p.ID),
			Name:  p.Name,
			Brand: p.Brand,
			Price: p.Price,
		}
	}

	return c.JSON(products)
}

// GetProductByID fetches a single product by ID using GORM
func GetProductByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var dbProduct ProductModel
	if err := db.First(&dbProduct, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
		}
		log.Printf("failed to fetch product by id: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch product"})
	}

	product := Product{
		ID:    int(dbProduct.ID),
		Name:  dbProduct.Name,
		Brand: dbProduct.Brand,
		Price: dbProduct.Price,
	}

	return c.JSON(product)
}

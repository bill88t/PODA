package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CorsMiddleware(port string) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "http://0.0.0.0:" + port,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	})
}

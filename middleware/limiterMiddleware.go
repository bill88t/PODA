package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// LimiterMiddleware rate limits
func LimiterMiddleware() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:               15000,
		Expiration:        15 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	})
}

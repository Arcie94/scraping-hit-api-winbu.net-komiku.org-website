package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimiter returns a rate limiting middleware
// Limits to 60 requests per minute per IP address
func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        60,              // Maximum 60 requests
		Expiration: 1 * time.Minute, // Per 1 minute window

		// Use IP address as unique identifier
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},

		// Custom response when limit is reached
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Rate limit exceeded",
				"message":     "You have exceeded the maximum of 60 requests per minute",
				"retry_after": "Please try again in 60 seconds",
			})
		},

		// Skip rate limiting for successful requests
		SkipFailedRequests: false,

		// Skip rate limiting for successful requests only (commented out - we want to limit all)
		// SkipSuccessfulRequests: false,
	})
}

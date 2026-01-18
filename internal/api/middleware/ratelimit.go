package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimiter creates a rate limiter middleware
// Default: 45 requests per minute per IP
func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        45,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Check for API key first
			apiKey := c.Get("X-API-Key")
			if apiKey != "" && IsValidAPIKey(apiKey) {
				return "apikey:" + apiKey
			}
			// Default to IP-based rate limiting
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "Too many requests. Please try again later.",
				},
			})
		},
		Storage: nil, // Use in-memory storage (default)
	})
}

// PremiumRateLimiter for API key users
// 450 requests per minute
func PremiumRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        450,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			apiKey := c.Get("X-API-Key")
			if apiKey != "" {
				return "premium:" + apiKey
			}
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "Premium rate limit exceeded.",
				},
			})
		},
	})
}

// IsValidAPIKey validates API keys
// TODO: Implement proper API key validation (database, env vars, etc.)
func IsValidAPIKey(key string) bool {
	// For now, simple validation
	// In production, check against database or environment variables
	validKeys := map[string]bool{
		"demo-key-12345": true,
		// Add more keys here or load from database
	}
	return validKeys[key]
}

// APIKeyMiddleware checks for valid API key and upgrades rate limit
func APIKeyMiddleware(c *fiber.Ctx) error {
	apiKey := c.Get("X-API-Key")

	if apiKey != "" && IsValidAPIKey(apiKey) {
		// Set flag for premium rate limit
		c.Locals("isPremium", true)
	}

	return c.Next()
}

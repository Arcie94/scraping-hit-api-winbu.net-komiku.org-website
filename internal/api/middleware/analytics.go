package middleware

import (
	"komiku-scraper/internal/analytics"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AnalyticsMiddleware tracks request analytics
func AnalyticsMiddleware(analyticsService *analytics.Analytics) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Track analytics
		if analyticsService != nil {
			event := analytics.Event{
				Timestamp:    start,
				Endpoint:     c.Path(),
				Method:       c.Method(),
				StatusCode:   c.Response().StatusCode(),
				ResponseTime: time.Since(start).Milliseconds(),
				IP:           c.IP(),
				APIKey:       c.Get("X-API-Key"),
				UserAgent:    c.Get("User-Agent"),
			}

			analyticsService.Track(event)
		}

		return err
	}
}

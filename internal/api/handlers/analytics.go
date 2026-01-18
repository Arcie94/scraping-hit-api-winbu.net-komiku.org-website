package handlers

import (
	"komiku-scraper/internal/analytics"
	"komiku-scraper/internal/api/middleware"
	"komiku-scraper/internal/models"

	"github.com/gofiber/fiber/v2"
)

// AnalyticsSummaryHandler returns analytics summary (protected)
func AnalyticsSummaryHandler(analyticsService *analytics.Analytics) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check API key
		apiKey := c.Get("X-API-Key")
		if !middleware.IsValidAPIKey(apiKey) {
			return c.Status(fiber.StatusUnauthorized).JSON(
				models.ErrorResponse("UNAUTHORIZED", "Valid API key required"),
			)
		}

		summary := analyticsService.GetSummary()
		return c.JSON(models.SuccessResponse(summary))
	}
}

// AnalyticsPopularHandler returns popular endpoints (protected)
func AnalyticsPopularHandler(analyticsService *analytics.Analytics) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check API key
		apiKey := c.Get("X-API-Key")
		if !middleware.IsValidAPIKey(apiKey) {
			return c.Status(fiber.StatusUnauthorized).JSON(
				models.ErrorResponse("UNAUTHORIZED", "Valid API key required"),
			)
		}

		popular, err := analyticsService.GetPopularEndpoints(10)
		if err != nil {
			return c.JSON(models.ErrorResponse("FETCH_FAILED", err.Error()))
		}

		return c.JSON(models.SuccessWithMeta(popular, len(popular)))
	}
}

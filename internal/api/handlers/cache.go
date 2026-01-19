package handlers

import (
	"komiku-scraper/scraper/cache"

	"github.com/gofiber/fiber/v2"
)

// GetCacheStats returns current cache statistics
func GetCacheStats(c *cache.Cache) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		stats := c.GetStats()
		return ctx.JSON(fiber.Map{
			"hits":     stats.Hits,
			"misses":   stats.Misses,
			"sets":     stats.Sets,
			"hit_rate": c.HitRate(),
			"size":     c.Size(),
		})
	}
}

// ClearCache removes all items from the cache
func ClearCache(c *cache.Cache) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		c.Clear()
		return ctx.JSON(fiber.Map{
			"message": "Cache cleared successfully",
		})
	}
}

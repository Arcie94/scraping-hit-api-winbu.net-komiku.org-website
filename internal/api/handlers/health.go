package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	startTime      time.Time
	requestCounter int64
	counterMutex   sync.Mutex
)

func init() {
	startTime = time.Now()
}

// IncrementRequestCounter increments the global request counter
func IncrementRequestCounter() {
	counterMutex.Lock()
	requestCounter++
	counterMutex.Unlock()
}

// HealthCheckHandler returns health status with scraper connectivity
func HealthCheckHandler(scraperURL string, scraperName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Test scraper connectivity
		scraperStatus := "ok"
		client := &http.Client{Timeout: 5 * time.Second}

		resp, err := client.Get(scraperURL)
		if err != nil || resp.StatusCode != 200 {
			scraperStatus = "error"
		}
		if resp != nil {
			resp.Body.Close()
		}

		// Calculate uptime
		uptime := time.Since(startTime).Seconds()

		return c.JSON(fiber.Map{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"uptime":    uptime,
			"scraper": fiber.Map{
				"name":   scraperName,
				"status": scraperStatus,
				"url":    scraperURL,
			},
			"requests_served": requestCounter,
		})
	}
}

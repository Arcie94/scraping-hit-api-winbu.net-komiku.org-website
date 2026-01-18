package handlers

import (
	"komiku-scraper/internal/models"
	"komiku-scraper/internal/service"

	"github.com/gofiber/fiber/v2"
)

// BatchAnimeRequest holds batch request data
type BatchAnimeRequest struct {
	Endpoints []string `json:"endpoints"`
}

// BatchAnimeHandler handles batch anime detail requests
func BatchAnimeHandler(svc *service.WinbuService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req BatchAnimeRequest
		if err := c.BodyParser(&req); err != nil {
			return c.JSON(models.ErrorResponse("INVALID_BODY", "Invalid request body"))
		}

		// Limit batch size
		if len(req.Endpoints) == 0 {
			return c.JSON(models.ErrorResponse("EMPTY_REQUEST", "Endpoints array is empty"))
		}
		if len(req.Endpoints) > 10 {
			return c.JSON(models.ErrorResponse("BATCH_TOO_LARGE", "Maximum 10 items per batch"))
		}

		// Fetch all anime details
		results := make([]interface{}, 0, len(req.Endpoints))
		errors := make([]fiber.Map, 0)

		for _, endpoint := range req.Endpoints {
			detail, err := svc.FetchAndParseDetail(endpoint)
			if err != nil {
				errors = append(errors, fiber.Map{
					"endpoint": endpoint,
					"error":    err.Error(),
				})
				continue
			}
			results = append(results, detail)
		}

		response := fiber.Map{
			"success": true,
			"data":    results,
			"meta": fiber.Map{
				"total":     len(results),
				"requested": len(req.Endpoints),
				"failed":    len(errors),
			},
		}

		if len(errors) > 0 {
			response["errors"] = errors
		}

		return c.JSON(response)
	}
}

// BatchMangaRequest holds batch manga request data
type BatchMangaRequest struct {
	Endpoints []string `json:"endpoints"`
}

// BatchMangaHandler handles batch manga detail requests
func BatchMangaHandler(svc *service.KomikuService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req BatchMangaRequest
		if err := c.BodyParser(&req); err != nil {
			return c.JSON(models.ErrorResponse("INVALID_BODY", "Invalid request body"))
		}

		// Limit batch size
		if len(req.Endpoints) == 0 {
			return c.JSON(models.ErrorResponse("EMPTY_REQUEST", "Endpoints array is empty"))
		}
		if len(req.Endpoints) > 10 {
			return c.JSON(models.ErrorResponse("BATCH_TOO_LARGE", "Maximum 10 items per batch"))
		}

		// Fetch all manga details
		results := make([]interface{}, 0, len(req.Endpoints))
		errors := make([]fiber.Map, 0)

		for _, endpoint := range req.Endpoints {
			detail, err := svc.FetchAndParseDetail(endpoint)
			if err != nil {
				errors = append(errors, fiber.Map{
					"endpoint": endpoint,
					"error":    err.Error(),
				})
				continue
			}
			results = append(results, detail)
		}

		response := fiber.Map{
			"success": true,
			"data":    results,
			"meta": fiber.Map{
				"total":     len(results),
				"requested": len(req.Endpoints),
				"failed":    len(errors),
			},
		}

		if len(errors) > 0 {
			response["errors"] = errors
		}

		return c.JSON(response)
	}
}

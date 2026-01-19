package handler

import (
	"komiku-scraper/internal/service"

	"github.com/gofiber/fiber/v2"
)

type WinbuHandler struct {
	Service *service.WinbuService
}

func NewWinbuHandler(svc *service.WinbuService) *WinbuHandler {
	return &WinbuHandler{Service: svc}
}

// Home Handler
func (h *WinbuHandler) Home(c *fiber.Ctx) error {
	data, err := h.Service.FetchHomeData()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	// Return the whole HomeData struct as JSON
	return c.JSON(data)
}

// Search Handler
func (h *WinbuHandler) Search(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Query parameter 'q' is required"})
	}
	results, err := h.Service.FetchSearch(query)
	if err != nil {
		// It might return empty list, not error, but handle error if any
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(results)
}

// Detail Handler (Anime Info)
func (h *WinbuHandler) Detail(c *fiber.Ctx) error {
	slug := c.Params("endpoint")
	// Winbu URL structure: https://winbu.net/anime/<slug>/ or /film/<slug>/ ?
	// Our service FetchAndParseDetail likely expects full URL.
	// CLI usually constructs it.
	// Let's assume standard anime slug for now. If it's a movie, the slug might be different path?
	// The CLI code (menu.go) handles "Search Anime" and "Film Terbaru" differently?
	// Actually menu.go: handleDetailWinbu(svc, scanner, targetAnime.Endpoint)
	// It passes the endpoint directly from the listing.
	// So if the user passes a slug via API, we need to know if it's movie or anime?
	// OR we just assume base URL + slug.
	// BUT, winbu endpoints in scraped lists are often full URLs or relative paths like https://winbu.net/anime/xyz/
	// If the user uses the search API, they get the full endpoint.
	// If they pass that endpoint to this API...
	// Let's assume the user passes the slug and we construct.
	// CAUTION: Winbu has /anime/slug and /film/slug.
	// Solution: Accept the FULL endpoint as query param? Or try both?
	// Better: The client should probably pass the full URL if possible, OR we make the API accept 'type' param.
	// For simplicity, let's look at how Komiku did it (constructed).
	// Winbu Service FetchAndParseDetail takes a full URL.
	// Let's try constructing /anime/ first.

	url := "https://winbu.net/anime/" + slug + "/"
	// Note: If it turns out to be a movie, this might fail 404.
	// Ideally we'd check or try both, but let's start with standard.
	// Or maybe the input 'endpoint' is actually the full URL encoded?
	// User request implication: "detail for film & drama".
	// Maybe we should allow 'type' param? /winbu/detail?slug=...&type=movie
	// Or simpler: /winbu/anime/:slug and /winbu/movie/:slug endpoints?

	data, err := h.Service.FetchAndParseDetail(url)
	if err != nil {
		// Try 'film' path if anime failed?
		url = "https://winbu.net/film/" + slug + "/"
		data, err = h.Service.FetchAndParseDetail(url)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	}
	return c.JSON(data)
}

// Episode Handler (Stream & Download)
func (h *WinbuHandler) Episode(c *fiber.Ctx) error {
	slug := c.Params("endpoint")
	// Episode URL: https://winbu.net/episode/<slug>/ (Verify this pattern)
	// Usually it's https://winbu.net/<slug>/ (direct) check parser.
	// Let's check ParseEpisodePage usage in CLI.

	url := "https://winbu.net/" + slug + "/"

	data, err := h.Service.FetchEpisode(url)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}

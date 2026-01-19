package handler

import (
	"fmt"
	"komiku-scraper/internal/service"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type KomikuHandler struct {
	Service *service.KomikuService
}

func NewKomikuHandler(svc *service.KomikuService) *KomikuHandler {
	return &KomikuHandler{Service: svc}
}

// Home Handler
func (h *KomikuHandler) Home(c *fiber.Ctx) error {
	data, err := h.Service.FetchHomeData()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}

// Search Handler
// Search Handler
func (h *KomikuHandler) Search(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Query parameter 'q' is required"})
	}

	// Construct Search URL (using main domain, not data subdomain)
	searchURL := fmt.Sprintf("https://komiku.id/?s=%s", strings.ReplaceAll(query, " ", "+"))

	results, err := h.Service.FetchAndParseList(searchURL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(results)
}

// Detail Handler
func (h *KomikuHandler) Detail(c *fiber.Ctx) error {
	slug := c.Params("endpoint")
	// Construct Full URL from slug
	url := "https://komiku.id/manga/" + slug + "/"

	data, err := h.Service.FetchAndParseDetail(url)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}

// Chapter Handler
func (h *KomikuHandler) Chapter(c *fiber.Ctx) error {
	slug := c.Params("endpoint")
	// Chapter URL: https://komiku.id/ch/<slug>/
	url := "https://komiku.id/ch/" + slug + "/"

	data, err := h.Service.FetchChapterImages(url)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}

// Genre Handler
func (h *KomikuHandler) Genres(c *fiber.Ctx) error {
	data, err := h.Service.FetchGenreList()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}

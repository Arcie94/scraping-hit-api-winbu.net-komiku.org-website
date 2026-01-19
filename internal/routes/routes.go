package routes

import (
	"komiku-scraper/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, komikuHandler *handler.KomikuHandler, winbuHandler *handler.WinbuHandler) {
	api := app.Group("/api/v1")

	// Komiku Routes
	komiku := api.Group("/komiku")
	komiku.Get("/home", komikuHandler.Home)
	komiku.Get("/search", komikuHandler.Search)
	komiku.Get("/manga/:endpoint", komikuHandler.Detail)
	komiku.Get("/chapter/:endpoint", komikuHandler.Chapter)
	komiku.Get("/genres", komikuHandler.Genres)

	// Winbu Routes
	winbu := api.Group("/winbu")
	winbu.Get("/home", winbuHandler.Home)
	winbu.Get("/search", winbuHandler.Search)
	winbu.Get("/detail/:endpoint", winbuHandler.Detail) // Handle both anime/movie via logic
	winbu.Get("/episode/:endpoint", winbuHandler.Episode)
}

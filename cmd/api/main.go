package main

import (
	"komiku-scraper/internal/handler"
	"komiku-scraper/internal/middleware"
	"komiku-scraper/internal/routes"
	"komiku-scraper/internal/service"
	"komiku-scraper/scraper/cache"
	"komiku-scraper/scraper/komiku"
	"komiku-scraper/scraper/winbu"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	log.Println("Starting Komiku & Winbu Scraper API...")

	// 1. Initialize Cache
	c := cache.New()

	// 2. Initialize HTTP Clients
	komikuClient := komiku.NewKomikuClient()
	winbuClient := winbu.NewWinbuClient()

	// 3. Initialize Services
	komikuService := service.NewKomikuService(komikuClient, c)
	winbuService := service.NewWinbuService(winbuClient, c)

	// 4. Initialize Handlers
	komikuHandler := handler.NewKomikuHandler(komikuService)
	winbuHandler := handler.NewWinbuHandler(winbuService)

	// 4. Initialize Fiber App
	app := fiber.New()
	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	app.Use(middleware.RateLimiter()) // Rate limiting: 60 req/min per IP

	// 5. Setup Routes
	routes.SetupRoutes(app, komikuHandler, winbuHandler)

	// 6. Start Server
	log.Fatal(app.Listen(":3000"))
}

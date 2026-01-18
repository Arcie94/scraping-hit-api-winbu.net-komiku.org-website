package main

import (
	"komiku-scraper/internal/service"
	"komiku-scraper/internal/ui"
	"komiku-scraper/scraper/komiku"
	"komiku-scraper/scraper/winbu"
)

func main() {
	// 1. Initialize HTTP Client (Scraper Core)
	komikuClient := komiku.NewKomikuClient()

	// 1b. Initialize Winbu Client
	winbuClient := winbu.NewWinbuClient()

	// 2. Initialize Service Layer (Business Logic)
	komikuSvc := service.NewKomikuService(komikuClient)
	winbuSvc := service.NewWinbuService(winbuClient)

	// 3. Start UI (Presentation Layer)
	ui.StartMenu(komikuSvc, winbuSvc)
}
